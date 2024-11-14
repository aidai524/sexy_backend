package http

import (
	"github.com/gin-gonic/gin"
	"sexy_backend/api/model"
	"sexy_backend/api/service"
	"sexy_backend/common/ecode"
	common "sexy_backend/common/http"
	dberror "sexy_backend/common/sexyerror"
)

// getAccountData 查询用户数据(super like数量相关)
// @Summary 查询用户数据(super like数量相关)
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} model.AccountData "data"
// @Router /account/data [get]
func getAccountData(c *gin.Context) {
	var (
		data *model.AccountData
		err  error
	)
	account := c.GetString("account")
	if account == "" {
		common.ReturnError(c, &dberror.Error{Code: ecode.RequestErr, Message: "account error"})
		return
	}

	data, err = service.API.GetAccountData(account)
	if err != nil {
		common.ReturnError(c, err)
		return
	}
	c.JSON(ecode.OK, common.Resp(ecode.OK, data))
}
