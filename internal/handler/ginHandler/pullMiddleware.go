package ginHandler

import (
	"ipfast_server/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

// 检查拉取次数
func PullMiddleware(c *gin.Context) {
	uidStr, exists := c.Get("user_id")
	if !exists {
		FailedResponse(c, "Authorization token is invalid")
		return
	}
	uid, ok := uidStr.(string)
	if !ok {
		FailedResponse(c, "Authorization token is invalid")
		return
	}
	err := services.CheckPullCount(uid, services.PullUpStrategy{
		Threshold: 3,
		Intervals: []time.Duration{
			time.Minute,
			5 * time.Minute,
			15 * time.Minute,
		},
	})
	if err != nil {
		FailedResponse(c, "pulling orders too frequently")
		return
	}
}
