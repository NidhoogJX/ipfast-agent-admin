package library

import (
	"fmt"
	"ipfast_server/internal/db/models"
	"ipfast_server/pkg/util/log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const ExpiresTime = 300

type VerificationCode = models.VerificationCode
type verificationCodeIP struct {
	IP        string
	ApplyTime int64
	Count     int
}

var verificationCodeIPList sync.Map

func init() {
	ClearVerificationCodeIPList()
}

func ClearVerificationCodeIPList() {
	verificationCodeIPList = sync.Map{}
}

/*
生成验证码

	param:
		email: 邮箱
		ip: IP 地址
	return:
		bool: 是否生成成功
*/
func MakeVerificationCode(email, ip string) (string, error) {
	value, ok := verificationCodeIPList.Load(ip)
	if ok {
		verificationCodeIP := value.(verificationCodeIP)
		if verificationCodeIP.Count >= 3 && verificationCodeIP.ApplyTime+60 > time.Now().Unix() {
			return "", fmt.Errorf("obtaining verification codes too frequently, please wait for one minute and try again")
		}
		verificationCodeIP.Count++
		verificationCodeIP.ApplyTime = time.Now().Unix()
		verificationCodeIPList.Store(ip, verificationCodeIP)
	} else {
		verificationCodeIPList.Store(ip, verificationCodeIP{IP: ip, ApplyTime: time.Now().Unix(), Count: 1})
	}
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	makecode := make([]byte, 6)
	for i := range makecode {
		makecode[i] = charset[rand.Intn(len(charset))]
	}
	model := &VerificationCode{}
	model.Code = string(makecode)
	model.ExpiresTime = ExpiresTime
	model.Email = email
	model.IP = ip
	err := model.CreateVerificationCode()
	if err != nil {
		return "", fmt.Errorf("failed to obtain verification code")
	}
	return model.Code, nil
}

/*
验证验证码

	param:
		email: 邮箱
		code: 验证码
	return:
		bool: 是否验证成功
*/
func VerifyVerificationCode(email, code string) bool {
	//TODO: 测试时使用 未接入邮箱时使用
	if code == "64a64a" {
		return true
	}
	model := &VerificationCode{}
	model.Code = code
	model.Email = email
	verificationCode, err := model.GetVerificationCode()
	if err != nil {
		return false
	}
	if !strings.EqualFold(verificationCode.Code, code) {
		return false
	}
	if verificationCode.CreatedTime+ExpiresTime < time.Now().Unix() {
		return false
	}
	err = model.UpdateVerificationCode()
	if err != nil {
		log.Error("update verification code error:%v", err)
	}
	return true
}
