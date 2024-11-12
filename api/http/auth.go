package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sexy_backend/api/conf"
	"sexy_backend/common/ecode"
	"sexy_backend/common/http"
	"sexy_backend/common/log"
	"sexy_backend/common/sexy/dao"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

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
		c.Next()
	}
}
