package services

import (
	"fmt"
	"ipfast_server/internal/db/models"
	"ipfast_server/pkg/util/log"
	"regexp"
	"time"

	"gorm.io/gorm"
)

type User = models.User

// GetUser 根据账号获取用户信息
func GetUserByAccount(accout string) (user *User, err error) {
	model := &User{}
	model.Email = accout
	model.Phone = accout
	user, err = model.FindByAccount()
	if err != nil {
		err = fmt.Errorf("failed to obtain user information")
	}
	return
}

// GetUser 根据账号获取用户信息
func GetUserByUserId(id int64) (user *User, err error) {
	model := &User{}
	model.Id = id
	user, err = model.FindById()
	if err != nil {
		err = fmt.Errorf("failed to obtain user information")
	}
	return
}

// 根据开始-截止时间和用户ID查询用户流量
func GetUserFlowByDate(startTime, endTime, uid int64) (userFlows []struct {
	Bytes int64  `json:"bytes"`
	Date  string `json:"date"`
}, err error) {
	return models.UserFlow{}.FindByDate(startTime, endTime, uid)
}

// 用户子账号分页查询
func GetSubUserByPage(page, pageSize int, uid int64, subUserName string) (results []models.SubUser, total int64, err error) {
	results, total, err = models.SubUser{}.GetSubUsers(page, pageSize, uid, subUserName)
	if err != nil {
		return nil, 0, fmt.Errorf("GetSubUserByPage failed to obtain subUser information")
	}
	return results, total, err
}

// 根据ID删除子用户
func DeleteSubUserByIDs(ids []int64, uid int64) error {
	model := &models.SubUser{}
	err := model.DeleteByIDs(ids, uid)
	if err != nil {
		return fmt.Errorf("failed to delete subUser")
	}
	// 删除子用户后，将需要删除的子用户ID列表添加到 kafka 消息队列中
	// data, err := json.Marshal(ids)
	// if err != nil {
	// 	fmt.Println("json marshal subUser failed", err)
	// 	return err
	// }
	// err = kafka.Producer.WriteMessage("delete", data)
	// if err != nil {
	// 	fmt.Println("send msg to kafka failed", err)
	// 	return err
	// }
	return nil
}

// 根据用户名模糊查询子用户信息
func SelectSubUserByUsername(username string, uid int64) (results models.SubUser, err error) {
	results, err = models.SubUser{}.GetSubUserByUsername(username, uid)
	if err != nil {
		return models.SubUser{}, fmt.Errorf("SelectSubUserByUsername failed to obtain subUser information")

	}
	return results, err
}

// 查询启用状态子用户信息
func SelectSubUserByEnableStatus(uid int64) (results []models.SubUser, err error) {
	results, err = models.SubUser{
		Status:       1,
		ParentUserID: uid,
	}.GetSubUserByStatus()
	if err != nil {
		err = fmt.Errorf("SelectSubUser failed to obtain subUser information")
	}
	return
}

// 添加子用户
func AddSubUser(uid int64, username, password string, maxCapacity float64, status, max_status int8, remarks string) (err error) {
	// 判断当前信息是否合法
	if len(username) < 1 || len(username) > 150 {
		return fmt.Errorf("username length error")
	}
	if len(password) < 1 || len(password) > 150 {
		return fmt.Errorf("password length error")
	}
	if maxCapacity < 0 {
		return fmt.Errorf("max capacity error")
	}
	if len(remarks) > 255 {
		return fmt.Errorf("remarks length error")
	}
	// 判断当前用户名是否已存在
	count, err := models.SubUser{}.IsUsernameExist(username, uid)
	if err != nil {
		return fmt.Errorf("failed to check username exist")
	}
	if count > 0 {
		return fmt.Errorf("username already exists")
	}

	model := &models.SubUser{
		ParentUserID: uid,
		Username:     username,
		Password:     password,
		MaxCapacity:  maxCapacity,
		Remarks:      remarks,
		Status:       status,
		CreatedTime:  time.Now().Unix(),
		UpdatedTime:  time.Now().Unix(),
		MaxStatus:    max_status,
	}
	_, err = model.Create()
	if err != nil {
		err = fmt.Errorf("registration failure")
		return err
	}
	// 新增子用户后，将其添加到 kafka 消息队列中
	// data, err := json.Marshal(subUser)
	// if err != nil {
	// 	fmt.Println("json marshal subUser failed", err)
	// 	return err
	// }
	// err = kafka.Producer.WriteMessage("insert", data)
	// if err != nil {
	// 	fmt.Println("send msg to kafka failed", err)
	// 	return err
	// }
	return err
}

// 更新子用户信息
func UpdateSubUser(uid, subUserId int64, username, password string, maxCapacity float64, remarks string, status, maxStatus int8) (err error) {
	// 判断当前信息是否合法
	if username != "" { //若用户名在参数中(进行修改)
		if len(username) < 1 || len(username) > 150 {
			return fmt.Errorf("username length error")
		}
	}
	if password != "" { //若密码在参数中(进行修改)
		if len(password) < 1 || len(password) > 150 {
			return fmt.Errorf("password length error")
		}
	}
	if status < 0 || status > 1 {
		return fmt.Errorf("status error")
	}
	if maxStatus < 0 || maxStatus > 1 {
		return fmt.Errorf("maxStatus error")
	}
	if maxCapacity < 0 {
		return fmt.Errorf("maxCapacity error")
	}
	if len(remarks) > 255 {
		return fmt.Errorf("remarks length error")
	}

	// 判断当前用户名是否已存在(需除去当前子账号)
	count, err := models.SubUser{}.IsUsernameExistExcludeCurrent(username, uid, subUserId)
	if err != nil {
		return fmt.Errorf("failed to check username exist")
	}
	if count > 0 {
		return fmt.Errorf("username already exists")
	}

	model := &models.SubUser{
		ParentUserID: uid,
		Username:     username,
		Password:     password,
		MaxCapacity:  maxCapacity,
		Status:       status,
		MaxStatus:    maxStatus,
		Remarks:      remarks,
		UpdatedTime:  time.Now().Unix(),
	}
	model.Id = subUserId
	model.ParentUserID = uid
	model.Username = username
	model.Password = password
	model.Remarks = remarks
	model.MaxCapacity = maxCapacity
	model.UpdatedTime = time.Now().Unix()

	err = model.Update()
	if err != nil {
		err = fmt.Errorf("failed to update subUser information")
		return err
	}
	// 修改子用户后，将其添加到 kafka 消息队列中
	// data, err := json.Marshal(model)
	// if err != nil {
	// 	fmt.Println("json marshal subUser failed", err)
	// 	return err
	// }
	// err = kafka.Producer.WriteMessage("update", data)
	// if err != nil {
	// 	fmt.Println("send msg to kafka failed", err)
	// 	return err
	// }
	return err
}

// 记录管理员登录信息
func RecordAdminLoginIpAndTime(admin models.Admin, ip string) error {
	admin.LoginIp = ip
	admin.LoginTime = time.Now().Unix()
	err := admin.UpdateLoginInfo()
	if err != nil {
		return fmt.Errorf("failed to record login information")
	}
	return nil
}

// 获取管理员信息
func GetAdminByUserId(uid int64) (admin models.Admin, err error) {
	admin, err = models.Admin{
		Id: uid,
	}.FindById()
	if err != nil {
		err = fmt.Errorf("failed to obtain admin information")
		return
	}
	return
}

// 获取代理商列表
func GetAgentList(page, size int, username string) (agentList []models.AgentInfo, total int64, err error) {
	agentList, total, err = models.Agent{}.SelectAgentList(page, size, username)
	if err != nil {
		err = fmt.Errorf("failed to obtain agent list")
	}
	return
}

// 添加代理商
func AddAgent(username, password, description string) (err error) {
	// 验证用户名、密码格式
	err = CheckNameAndPassword(username, password)
	if err != nil {
		return err
	}
	// 验证代理商用户名是否重复
	count, err := models.Agent{
		Name: username,
	}.IsExistByName()
	if err != nil {
		return fmt.Errorf("failed to check agent exist")
	}
	if count != 0 {
		return fmt.Errorf("username already exists")
	}
	salt, err := generateSalt(6)
	if err != nil {
		return fmt.Errorf("failed to generate salt")
	}
	now := time.Now().Unix()
	err = models.Agent{
		Name:        username,
		Password:    generateMD5(password + salt),
		AppKey:      generateMD5(username + salt),
		Salt:        salt,
		Description: description,
		Status:      1,
		CreateTime:  now,
		UpdateTime:  now,
	}.Create()
	if err != nil {
		return fmt.Errorf("failed to add agent")
	}
	return
}

// 编辑代理商信息
func EditAgent(agentId int64, password, description string, status int8) (err error) {
	agent, err := models.Agent{
		Id: agentId,
	}.FindById()
	if err != nil {
		err = fmt.Errorf("failed to obtain agent information")
		return
	}
	if password != agent.Password {
		// 检查密码格式，(跳过用户名)
		err = CheckNameAndPassword("username", password)
		if err != nil {
			return err
		}
		salt, err := generateSalt(6)
		if err != nil {
			err = fmt.Errorf("failed to generate salt")
			return err
		}
		agent.Password = generateMD5(password + salt)
		agent.Salt = salt
	}
	now := time.Now().Unix()
	agent.Status = status
	agent.Description = description
	agent.UpdateTime = now
	agent.UpdateInfo()
	if err != nil {
		err = fmt.Errorf("failed to update user information")
	}
	return
}

// 给代理商充值流量
func RechargeFlowToAgent(agentId, count int64, sign, description string) (err error) {
	// 检查代理商是否存在
	agent, err := models.Agent{
		Id: agentId,
	}.FindById()
	if err != nil {
		return fmt.Errorf("failed to obtain agent information")
	}
	if agent.Id == 0 {
		return fmt.Errorf("agent not found")
	}
	if agent.Status != 1 {
		return fmt.Errorf("agent is disabled")
	}
	nowTime := time.Now().Unix()
	now := time.Now()
	msec := now.UnixMilli() % 1000
	rechargeId := fmt.Sprintf("%04d%02d%02d%02d%02d%02d%03d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), msec)
	flow := count * 1024 * 1024 * 1024
	if sign == "+" {
		if count > 100000000 {
			return fmt.Errorf("recharge count is too large")
		}
		// 开启事务处理
		err = models.DB.Instance.Transaction(func(tx *gorm.DB) error {
			// 为代理商添加流量
			err = models.Agent{
				Id: agentId,
			}.RechargeFlow(flow, nowTime)
			if err != nil {
				log.Error("failed to recharge flow to agent, agentId: %d", agentId)
				return fmt.Errorf("failed to recharge flow to agent")
			}
			// 添加流量充值记录
			err = models.Recharge{
				Id:          rechargeId,
				AgentID:     agentId,
				Count:       count,
				Unit:        1,
				PayMethod:   "管理员充值",
				Description: description,
				Status:      1,
				CreateTime:  nowTime,
				UpdateTime:  nowTime,
			}.Create()
			if err != nil {
				log.Error("failed to add recharge record to agent, agentId: %d", agentId)
				return fmt.Errorf("failed to add recharge record to agent")
			}
			return nil
		})
	} else {
		if agent.TotalFlow < flow {
			log.Error("reduced flow more than total flow of agents, agentId: %d", agentId)
			return fmt.Errorf("reduced flow more than total flow of agents")
		}
		// 开启事务处理
		err = models.DB.Instance.Transaction(func(tx *gorm.DB) error {
			// 减少代理商流量
			err = models.Agent{
				Id: agentId,
			}.ReduceFlow(flow, nowTime)
			if err != nil {
				log.Error("failed to reduce flow to agent, agentId: %d", agentId)
				return fmt.Errorf("failed to reduce flow to agent")
			}
			// 添加流量减少记录
			err = models.Recharge{
				Id:          rechargeId,
				AgentID:     agentId,
				Count:       -count,
				Unit:        1,
				PayMethod:   "管理员充值",
				Description: description,
				Status:      1,
				CreateTime:  nowTime,
				UpdateTime:  nowTime,
			}.Create()
			if err != nil {
				log.Error("failed to add recharge record to agent, agentId: %d", agentId)
				return fmt.Errorf("failed to add recharge record to agent")
			}
			return nil
		})
	}
	if err != nil {
		log.Info("分配流量成功, agentId: %d, sign: %s, count: %d", agentId, sign, count)
	}
	return
}

// 检查用户名和密码
func CheckNameAndPassword(customerName, password string) error {
	if len(customerName) < 4 || len(customerName) > 20 {
		return fmt.Errorf("username length must be between 4 and 20")
	}
	if len(password) < 8 || len(password) > 16 {
		return fmt.Errorf("password length must be between 8 and 20")
	}
	if IsValidEmail(customerName) {
		return fmt.Errorf("username cannot be in email format")
	}
	// 验证密码格式
	re := regexp.MustCompile(`^[A-Za-z0-9!@#$%^&\-\*()_+\]\[\}\{|;:,.<>?]+$`)
	if re.MatchString(customerName) {
		if re.MatchString(password) {
			return nil
		} else {
			return fmt.Errorf("password can only be composed of english, numbers, or special symbols")
		}
	} else {
		return fmt.Errorf("username can only be composed of english, numbers, or special symbols")
	}
}

// 获取代理商的用户列表
func GetUserList(page, size int, agentId int64, username string) (userList []models.UserInfo, total int64, err error) {
	userList, total, err = models.User{}.SelectUserList(page, size, agentId, username)
	if err != nil {
		err = fmt.Errorf("failed to obtain user list")
	}
	return
}
