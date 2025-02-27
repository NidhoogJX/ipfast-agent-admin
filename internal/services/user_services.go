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

// 获取代理商信息
func GetAgentInfo(uid int64) (agent models.Agent, err error) {
	agent, err = models.Agent{
		Id: uid,
	}.FindById()
	if err != nil {
		err = fmt.Errorf("failed to obtain agent information")
	}
	return
}

// 记录代理商登录信息
func RecordAgentLoginIpAndTime(agent models.Agent, ip string) error {
	agent.LoginIp = ip
	agent.LoginTime = time.Now().Unix()
	err := agent.UpdateLoginInfo()
	if err != nil {
		return fmt.Errorf("failed to record login information")
	}
	return nil
}

// 获取代理商的流量统计
func GetTotalFlowDetail(id int64) (flowDetail models.AgentFlowInfo, err error) {
	flowDetail, err = models.Agent{}.SelectTotalFlowInfo(id)
	if err != nil {
		err = fmt.Errorf("failed to obtain flow detail")
	}
	return
}

// 获取代理商当天的流量统计
func GetCurrentFlowDetail(id int64) (currentFlowDetail models.CurrentAgentFlowInfo, err error) {
	// 获取当日0点时间戳
	t := time.Now()
	currentTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	startTime := currentTime.Unix()
	endTime := startTime + 86400
	currentFlowDetail, err = models.Agent{}.SelectCurrentFlowInfo(id, startTime, endTime)
	if err != nil {
		err = fmt.Errorf("failed to obtain user current flow detail")
	}
	return
}

// 获取用户的流量明细
func GetUserFlowDetail(id, startTime, endTime int64) (flowDate []models.DateFlow, err error) {
	flowDate, err = models.Agent{}.SelectUserFlowInfo(id, startTime, endTime)
	if err != nil {
		err = fmt.Errorf("failed to obtain user flow detail")
	}
	return
}

// 用户列表查询
func GetUserListByPage(agentId int64, username string, page, pageSize int, status, totalSort, usedSort, enableSort int8) (userList []models.UserInfo, total int64, err error) {
	userList, total, err = models.User{}.SelectUserList(agentId, username, page, pageSize, status, totalSort, usedSort, enableSort)
	if err != nil {
		return []models.UserInfo{}, total, fmt.Errorf("failed to obtain user list")
	}
	return
}

// 代理商添加用户
func AddUser(agentId int64, username, password, description string) (err error) {
	// 检查用户名、密码格式
	err = CheckNameAndPassword(username, password)
	if err != nil {
		return err
	}
	// 检查用户名是否重复
	count, err := models.User{
		Name: username,
	}.CheckUsernameExist()
	if err != nil {
		err = fmt.Errorf("failed to check username exist")
		return err
	}
	if count > 0 {
		return fmt.Errorf("username already exists")
	}
	salt, err := generateSalt(6)
	if err != nil {
		err = fmt.Errorf("failed to generate salt")
	}
	now := time.Now().Unix()
	models.User{
		Name:        username,
		AgentId:     agentId,
		Password:    generateMD5(password + salt),
		Salt:        salt,
		AppKey:      generateMD5(username + salt),
		Description: description,
		Status:      1,
		CreateTime:  now,
		UpdateTime:  now,
	}.Create()
	return
}

// 代理商修改用户信息
func EditUser(userId int64, password, description string, status int8) (err error) {
	user, err := models.User{
		Id: userId,
	}.SelectById()
	if err != nil {
		err = fmt.Errorf("failed to obtain user information")
		return
	}
	if password != user.Password {
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
		user.Password = generateMD5(password + salt)
		user.Salt = salt
	}
	now := time.Now().Unix()
	err = models.User{
		Id:          userId,
		Password:    user.Password,
		Salt:        user.Salt,
		Description: description,
		Status:      1,
		UpdateTime:  now,
	}.UpdateUserInfo()
	if err != nil {
		err = fmt.Errorf("failed to update user information")
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

// 代理商为用户分配流量
func DistributeFlowToUser(agentId, userId int64, count float64) (err error) {
	now := time.Now().Unix()
	flow := int64(count * 1024 * 1024 * 1024)
	// 开启事务处理
	err = models.DB.Instance.Transaction(func(tx *gorm.DB) error {
		// 判断代理商流量余额是否足够
		agent, err := models.Agent{
			Id: agentId,
		}.FindById()
		if err != nil {
			log.Error("failed to obtain agent information, agentId: %d", agentId)
			return fmt.Errorf("failed to obtain agent information")
		}
		if agent.TotalFlow-agent.DistributeFlow < flow {
			return fmt.Errorf("not enough flow")
		}
		// 检查用户是否存在,且可正常使用
		user, err := models.User{
			Id: userId,
		}.FindById()
		if err != nil || user.Id == 0 {
			return fmt.Errorf("failed to obtain user information")
		}
		if user.Status != 1 {
			return fmt.Errorf("user is disabled")
		}
		// 为用户添加流量
		err = models.FlowRecord{
			UserID:        userId,
			AgentId:       agentId,
			Type:          2, // 流量记录类型(1:购买流量/管理员分配,2:代理商分配流量)
			PurchasedFlow: flow,
			Deadline:      now + 86400*365,
			OrderId:       "1", // 订单ID(管理员添加的流量为0;代理商分配的为1)
			CreatedTime:   now,
			UpdatedTime:   now,
		}.Create()
		if err != nil {
			log.Error("failed to add flow record to user, userId: %d, agentId: %d", userId, agentId)
			return fmt.Errorf("failed to add flow record to user")
		}
		// 扣除代理商流量
		agent.DistributeFlow += flow
		agent.UpdateTime = now
		err = agent.UpdateFlowInfo()
		if err != nil {
			log.Error("failed to update agent flow information, agentId: %d", agentId)
			return fmt.Errorf("failed to update agent flow information")
		}
		return nil
	})
	if err != nil {
		log.Info("分配流量成功,agentId: %d, userId: %d, flow: %d", agentId, userId, flow)
	}
	return
}

// 代理商的流量分配日志
func GetDistributeFlowLog(agentId int64, page, pageSize int, username string) (distributeLogList []models.DistributeFlowLog, total int64, err error) {
	distributeLogList, total, err = models.FlowRecord{}.SelectDistributeFlowLog(agentId, page, pageSize, username)
	if err != nil {
		err = fmt.Errorf("failed to obtain distribute flow log")
	}
	return
}
