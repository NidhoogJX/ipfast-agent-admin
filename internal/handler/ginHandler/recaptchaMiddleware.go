package ginHandler

import (
	"encoding/json"
	"ipfast_server/internal/handler/network/request"
	"ipfast_server/pkg/util/log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const recaptchaSecret = "0x4AAAAAAA4eeKgtHgnMlcJeVd3mp1emw6k"

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
	Cdata       string   `json:"cdata"`
}

/*
验证 reCAPTCHA
响应示例:

	{
	  "success": true,
	  "challenge_ts": "2022-02-28T15:14:30.096Z",
	  "hostname": "example.com",
	  "error-codes": [],
	  "action": "login",
	  "cdata": "sessionid-123456789"
	}
*/
func VerifyRecaptcha(c *gin.Context) {
	cfToken := c.GetHeader("cf-turnstile")
	if cfToken == "" {
		FailedResponse(c, "recaptcha is invalid")
		return
	}
	reponse := &struct {
		Success     bool   `json:"success"`
		ChallengeTS string `json:"challenge_ts"`
		Hostname    string `json:"hostname"`
		Action      string `json:"action"`
	}{}
	body, _, err := request.Post("https://challenges.cloudflare.com/turnstile/v0/siteverify",
		map[string]string{
			"secret":   recaptchaSecret,
			"response": cfToken,
			"remoteip": c.ClientIP(),
		})
	if err != nil {
		FailedResponse(c, "verify recaptcha is error")
		return
	}
	log.Debug("recaptcha response:%s", string(body))
	err = json.Unmarshal(body, reponse)
	log.Info("recaptcha response:%v", reponse)
	if err != nil || !reponse.Success || reponse.Action != "login" || strings.Contains(reponse.Hostname, "ipfast.com") || reponse.ChallengeTS == "" {
		log.Error("recaptcha response error:%s,response%v ", err, strings.Contains(reponse.Hostname, "ipfast.com"))
		FailedResponse(c, "verify recaptcha is error")
		return
	}
	challengeTime, err := time.Parse(time.RFC3339, reponse.ChallengeTS)
	if err != nil {
		log.Error("parse challenge time error:%s", err)
		FailedResponse(c, "verify recaptcha is error")
		return
	}
	// 允许的时间范围，例如 2 分钟
	allowedDuration := 2 * time.Minute
	if time.Since(challengeTime) > allowedDuration {
		log.Error("challenge time is expired")
		FailedResponse(c, "verify recaptcha is error")
		return
	}
	c.Next()
}
