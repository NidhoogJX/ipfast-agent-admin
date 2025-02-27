package services

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"regexp"
)

// GenerateMD5 生成字符串的 MD5 哈希值
func generateMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// GenerateMD5 生成字符串的 MD5 哈希值
func GenerateMD5(text string) string {
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

// GenerateSalt 生成一个随机盐
func GenerateSalt(size int) (string, error) {
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
// func isValidPhone(phone string) bool {
// 	if len(phone) != 11 {
// 		return false
// 	}
// 	re := regexp.MustCompile(`^1[3456789]\d{9}$`)
// 	return re.MatchString(phone)
// }

// isValidPassword 验证密码长度
func isValidPassword(password string) bool {
	return len(password) >= 8 && len(password) <= 20
}
