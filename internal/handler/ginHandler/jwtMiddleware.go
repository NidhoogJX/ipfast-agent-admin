package ginHandler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("7ghsjh*&*^&^%$#")

// 生成 JWT 令牌
func GenerateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userId,
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(15 * 24 * time.Hour).Unix(),
	})
	return token.SignedString(jwtKey)
}

// 校验 JWT 令牌
func ValidateJWT(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		FailedResponseCode(c, 2, "Authorization is invalid")
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		FailedResponseCode(c, 2, "Authorization is invalid")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		FailedResponseCode(c, 2, "Authorization is invalid")
		return
	}

	// 验证过期时间
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			FailedResponseCode(c, 2, "Authorization is expired")
			return
		}
	} else {
		FailedResponseCode(c, 2, "Authorization is invalid")
		return
	}

	// 验证生效时间
	if nbf, ok := claims["nbf"].(float64); ok {
		if time.Unix(int64(nbf), 0).After(time.Now()) {
			FailedResponseCode(c, 2, "Authorization is not yet valid")
			return
		}
	} else {
		FailedResponseCode(c, 2, "Authorization is invalid")
		return
	}

	c.Set("user_id", claims["uid"])
	c.Next()
}
