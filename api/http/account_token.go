package http

import (
	"github.com/gin-gonic/gin"
	"sexy_backend/api/service"
	"sexy_backend/common/ecode"
	gin2 "sexy_backend/common/gin"
	"sexy_backend/common/http"
	dberror "sexy_backend/common/sexyerror"
)

type AccountTokenParam struct {
	Address   string `form:"address" json:"address" binding:"required"`
	Signature string `form:"signature" json:"signature" binding:"required"`
	Time      int64  `form:"time" json:"time" binding:"required"`
}

// getAccountToken 获取用户token
// @Summary 获取用户token
// @Tags 全局
// @Accept application/json
// @Produce application/json
// @Param object query AccountTokenParam false "查询参数"
// @Success 200 {string} string "返回的用户 token"
// @Router /account/token [get]
func getAccountToken(c *gin.Context) {
	var (
		param AccountTokenParam
		token string
		err   error
	)
	if err = gin2.ShouldBind(c, &param); err != nil {
		common.ReturnError(c, &dberror.Error{Code: ecode.RequestErr, Message: err.Error()})
		return
	}

	token, err = service.API.GetAccountToken(param.Address, param.Time, param.Signature)
	if err != nil {
		common.ReturnError(c, err)
		return
	}
	c.JSON(ecode.OK, common.Resp(ecode.OK, token))
}
