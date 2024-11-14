package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sexy_backend/api/conf"
	"sexy_backend/common/ecode"
	"sexy_backend/common/http"
	"sexy_backend/common/log"
	"sexy_backend/common/sexy/dao"
	"sexy_backend/common/sexy/database"
	"time"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		authFunc(token, c)
	}
}

func authFunc(token string, c *gin.Context) {
	account, err := dao.GetAuth(conf.Conf.Red, token)
	if err != nil {
		log.Error("GetAuth err: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.RespWithMsg(ecode.ServerErr, "error"))
		return
	}
	if len(account) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.RespWithMsg(ecode.ServerErr, "token error"))
		return
	}
	c.Set("account", account)

	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	err = database.CreateAccountActive(conf.Conf.DB, account, t.UnixMilli())
	if err != nil {
		log.Error("CreateAccountActive err: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.RespWithMsg(ecode.ServerErr, "create account active error"))
		return
	}
	c.Next()
}

// checkAuthorization 判断是否需要认证的中间件
func checkAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的 Authorization
		token := c.GetHeader("Authorization")

		// 如果 Authorization 头存在并且是有效的
		if token != "" {
			authFunc(token, c)
		} else {
			c.Next()
		}
	}
}
