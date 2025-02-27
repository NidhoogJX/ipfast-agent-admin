package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var jwtKey = []byte("7ghsjh*&*^&^%$#")

// 生成 JWT 令牌
func GenerateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userId,
		"nbf": 1704052173,
		"exp": 1704252173,
	})
	return token.SignedString(jwtKey)
}

// 校验 JWT 令牌
func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("authorization is invalid: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("authorization is invalid")
	}

	// 验证过期时间
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return "", fmt.Errorf("authorization is expired")
		}
	} else {
		return "", fmt.Errorf("authorization is invalid")
	}

	// 验证生效时间
	if nbf, ok := claims["nbf"].(float64); ok {
		if time.Unix(int64(nbf), 0).After(time.Now()) {
			return "", fmt.Errorf("authorization is not yet valid")
		}
	} else {
		return "", fmt.Errorf("authorization is invalid")
	}

	return claims["uid"].(string), nil
}

// 测试生成 JWT 令牌
func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken("25")
	fmt.Println(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

// 测试校验 JWT 令牌
func TestValidateToken(t *testing.T) {
	// 生成一个有效的 JWT 令牌
	token, err := GenerateToken("25")
	assert.NoError(t, err)

	// 校验 JWT 令牌
	userId, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "12345", userId)
}

// 测试过期的 JWT 令牌
func TestExpiredToken(t *testing.T) {
	// 生成一个过期的 JWT 令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": "12345",
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(-time.Hour).Unix(), // 设置过期时间为过去的时间
	})
	expiredToken, err := token.SignedString(jwtKey)
	assert.NoError(t, err)

	// 校验过期的 JWT 令牌
	_, err = ValidateToken(expiredToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authorization is expired")
}

// 测试无效的 JWT 令牌
func TestInvalidToken(t *testing.T) {
	// 创建一个无效的 JWT 令牌
	invalidToken := "invalid.token.here"

	// 校验无效的 JWT 令牌
	_, err := ValidateToken(invalidToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authorization is invalid")
}
