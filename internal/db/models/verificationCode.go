package models

import (
	"time"
)

type VerificationCode struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Code        string `gorm:"not null" json:"code"`
	Email       string `gorm:"not null;index" json:"email"`
	IP          string `json:"ip" gorm:"default:'';index"`
	CreatedTime int64  `json:"created_time"`
	ExpiresTime int64  `gorm:"default:0;index" json:"expires_time"`
}

func (VerificationCode) TableName() string {
	return "ip_verification_codes"
}

// 查询邮箱验证码
func (model *VerificationCode) GetVerificationCode() (verificationCode *VerificationCode, err error) {
	if err = DB.Where("email = ? and expires_time != 0", model.Email).Order("created_time desc").
		First(&verificationCode).Error; err != nil {
		return nil, err
	}
	return
}

// 新增验证码
func (model *VerificationCode) CreateVerificationCode() error {
	model.CreatedTime = time.Now().Unix()
	if err := DB.Create(model).Error; err != nil {
		return err
	}
	return nil
}

// 使用验证码
func (model *VerificationCode) UpdateVerificationCode() error {
	if err := DB.Table(model.TableName()).Where("email = ? and expires_time != 0", model.Email).Update("expires_time", 0).Error; err != nil {
		return err
	}
	return nil
}
