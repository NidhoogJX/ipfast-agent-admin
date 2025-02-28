package library

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"ipfast_server/internal/db/models"
	"ipfast_server/pkg/util/log"
	"regexp"
	"time"
)

type User = models.User

// GenerateMD5 生成字符串的 MD5 哈希值
func generateMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// GenerateSalt 生成一个随机盐
func generateSalt(size int) (string, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// isValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// isValidPhone 验证手机号格式
func isValidPhone(phone string) bool {
	if len(phone) != 11 {
		return false
	}
	re := regexp.MustCompile(`^1[3456789]\d{9}$`)
	return re.MatchString(phone)
}

// isValidPassword 验证密码长度
func isValidPassword(password string) bool {
	return len(password) >= 8 && len(password) <= 20
}

// GetUser 根据账号获取用户信息
func GetUser(accout string) (user *User, err error) {
	model := &User{}
	model.Email = accout
	model.Phone = accout
	user, err = model.FindByAccount()
	if err != nil {
		err = fmt.Errorf("failed to obtain user information")
	}
	return
}

/*
管理员账号密码登陆
*/
func Auth(name, password string) (admin models.Admin, err error) {
	admin, err = models.Admin{
		Name: name,
	}.FindByName()
	if err != nil || admin.Id == 0 {
		err = fmt.Errorf("account does not exist/has been cancelled")
		return
	}
	if admin.Password != generateMD5(password+admin.Salt) {
		err = fmt.Errorf("account or password error")
		return
	}
	if admin.Status == 0 {
		err = fmt.Errorf("account has been disabled")
		return
	}
	return
}

/*
用户注册
*/
func Register(email, phone, password string) (err error) {
	if (email == "" && phone == "") || password == "" {
		return fmt.Errorf("param error")
	}
	if email != "" && !IsValidEmail(email) {
		return fmt.Errorf("email format error")
	}
	if phone != "" && !isValidPhone(phone) {
		return fmt.Errorf("phone format error")
	}
	if !isValidPassword(password) {
		return fmt.Errorf("password length must be between 8 and 20")
	}
	salt, err := generateSalt(6)
	if err != nil {
		return err
	}
	timeUnix := time.Now().Unix()
	model := &models.User{
		Email:      email,
		Phone:      phone,
		Password:   generateMD5(password + salt),
		AppKey:     generateMD5(email + phone + salt),
		Salt:       salt,
		CreateTime: timeUnix, // 注册时间
		UpdateTime: timeUnix, // 更新时间
		Status:     1,        // 正常状态
		Money:      0,
	}

	user, err := model.FindByAccount()
	if err != nil {
		log.Error("查询用户失败:%v", err)
		return fmt.Errorf("check if the username registration failed")
	}
	if user.Id > 0 {
		return fmt.Errorf("the username or email already exists, please change it")
	}
	err = model.Create()
	if err != nil {
		log.Error("创建新用户失败:%v", err)
		err = fmt.Errorf("registration failure")
	}
	return
}

/*
用户登出
*/
func Logout(userID int64) error {
	return nil
}

/*
重置密码
*/
func ResetPassword(uid int64, password string) error {
	if !isValidPassword(password) {
		return fmt.Errorf("password length must be between 8 and 20")
	}
	salt, err := generateSalt(6)
	if err != nil {
		return err
	}
	model := &User{
		Id:       uid,
		Password: generateMD5(password + salt),
		Salt:     salt,
	}
	err = model.UpdatePassword()
	if err != nil {
		return fmt.Errorf("failed to reset password")
	}
	return nil
}

// GetUser 根据账号获取用户信息
func GetUserById(id int64) (user *User, err error) {
	model := &User{}
	model.Id = id
	user, err = model.FindById()
	if err != nil {
		err = fmt.Errorf("failed to obtain user information")
	}
	return
}

// func GetApiProxyUrl() string {
// 	param := & struct{
// 		Number int `json:"num" binding:"required;" `
// 		"num": 5,
// 		"countrty":"",
// 		"provice":"",
// 		"city":"",
// 		"protocol": "http",
// 		"return_type": "txt",
// 		"ib": 0,
// 		"separator":"\t"

// 	}
// 	return "http://
// }
