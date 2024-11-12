package http

import (
	"github.com/gin-gonic/gin"
	"sexy_backend/api/service"
	"sexy_backend/common/ecode"
	gin2 "sexy_backend/common/gin"
	common "sexy_backend/common/http"
	"sexy_backend/common/sexy/database"
	dberror "sexy_backend/common/sexyerror"
)

type PostProjectParam struct {
	TokenName         string `json:"token_name" form:"token_name" binding:"required,max=100"`                  // 项目代币名称(唯一)
	TokenSymbol       string `json:"token_symbol" form:"token_symbol" binding:"required,max=50"`               // 项目代币symbol(唯一)
	Icon              string `json:"icon" form:"icon" binding:"required,max=300"`                              // 项目图标URL
	Video             string `json:"video" form:"video" binding:"max=300"`                                     // 视频URL
	BackgroundStory   string `json:"background_story" form:"background_story" binding:"required,max=2000"`     // 背景故事
	FutureDevelopment string `json:"future_development" form:"future_development" binding:"required,max=2000"` // 未来发展
	WhitePaper        string `json:"white_paper" form:"white_paper" binding:"required,max=300"`                // 白皮书URL
	X                 string `json:"x" form:"x" binding:"max=300"`                                             // 项目方的社交媒体信息
	Tg                string `json:"tg" form:"tg" binding:"max=300"`                                           // 项目方的社交媒体信息
	Country           string `json:"country" form:"country" binding:"required,max=300"`                        // 国家信息
}

// postProject 上传项目
// @Summary 上传项目
// @Tags 项目
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param object query PostProjectParam false "查询参数"
// @Router /project [post]
func postProject(c *gin.Context) {
	var (
		param PostProjectParam
		err   error
	)

	if err = gin2.ShouldBind(c, &param); err != nil {
		common.ReturnError(c, &dberror.Error{Code: ecode.RequestErr, Message: err.Error()})
		return
	}
	address := c.GetString("account")
	if address == "" {
		common.ReturnError(c, &dberror.Error{Code: ecode.RequestErr, Message: "account error"})
		return
	}

	err = service.API.PostProject(address, param.TokenName, param.TokenSymbol, param.Icon, param.Video, param.BackgroundStory, param.FutureDevelopment, param.WhitePaper, param.X, param.Tg, param.Country)
	if err != nil {
		common.ReturnError(c, err)
		return
	}
	c.JSON(ecode.OK, common.Resp(ecode.OK, nil))
}

type GetProjectParam struct {
	ID          uint64 `json:"id" form:"id"`
	TokenName   string `json:"token_name" form:"token_name" binding:"max=100"`    // 项目代币名称(唯一)
	TokenSymbol string `json:"token_symbol" form:"token_symbol" binding:"max=50"` // 项目代币symbol(唯一)
}

// getProject 查询项目
// @Summary 查询项目
// @Tags 项目
// @Accept application/json
// @Produce application/json
// @Param object query GetProjectParam false "查询参数"
// @Success 200 {object} database.Project "data"
// @Router /project [get]
func getProject(c *gin.Context) {
	var (
		param   GetProjectParam
		project *database.Project
		err     error
	)

	if err = gin2.ShouldBind(c, &param); err != nil {
		common.ReturnError(c, &dberror.Error{Code: ecode.RequestErr, Message: err.Error()})
		return
	}
	project, err = service.API.GetProject(param.ID, param.TokenName, param.TokenSymbol)
	if err != nil {
		common.ReturnError(c, err)
		return
	}
	c.JSON(ecode.OK, common.Resp(ecode.OK, project))
}

type GetProjectSearchParam struct {
	Text   string `json:"text" form:"text" binding:"required,max=100"`
	Limit  int    `json:"limit" form:"limit" binding:"required,max=100"`
	Offset int    `json:"offset" form:"offset"`
}

// getProject 模糊搜索项目
// @Summary 模糊搜索项目
// @Tags 项目
// @Accept application/json
// @Produce application/json
// @Param object query GetProjectSearchParam false "查询参数"
// @Success 200 {array} database.Project ""
// @Router /project/search [get]
func getProjectSearch(c *gin.Context) {
	var (
		param       GetProjectSearchParam
		projectList []*database.Project
		haxNextPage bool
		err         error
	)

	if err = gin2.ShouldBind(c, &param); err != nil {
		common.ReturnError(c, &dberror.Error{Code: ecode.RequestErr, Message: err.Error()})
		return
	}
	projectList, haxNextPage, err = service.API.GetProjectSearch(param.Text, param.Limit, param.Offset)
	if err != nil {
		common.ReturnError(c, err)
		return
	}
	c.JSON(ecode.OK, common.Resp(ecode.OK, common.NewHasNextPageResult{Data: projectList, HasNextPage: haxNextPage}))
}

type GetProjectListParam struct {
	Limit int `json:"limit" form:"limit" binding:"required,max=100"`
}

// getProject 项目列表
// @Summary 项目列表
// @Tags 项目
// @Accept application/json
// @Produce application/json
// @Param object query GetProjectListParam false "查询参数"
// @Success 200 {array} database.Project ""
// @Router /project/list [get]
func getProjectList(c *gin.Context) {
	var (
		param       GetProjectListParam
		projectList []*database.Project
		haxNextPage bool
		err         error
	)

	if err = gin2.ShouldBind(c, &param); err != nil {
		common.ReturnError(c, &dberror.Error{Code: ecode.RequestErr, Message: err.Error()})
		return
	}
	account := c.GetString("account")
	if account == "" {
		common.ReturnError(c, &dberror.Error{Code: ecode.RequestErr, Message: "account error"})
		return
	}

	projectList, haxNextPage, err = service.API.GetProjectList(account, param.Limit)
	if err != nil {
		common.ReturnError(c, err)
		return
	}
	c.JSON(ecode.OK, common.Resp(ecode.OK, common.NewHasNextPageResult{Data: projectList, HasNextPage: haxNextPage}))
}

type PostProjectLikeParam struct {
	ID uint64 `json:"id" form:"id" binding:"required"`
}

// postProjectLike like项目
// @Summary like项目
// @Tags 项目
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param object query PostProjectLikeParam false "查询参数"
// @Router /project/like [post]
func postProjectLike(c *gin.Context) {
	var (
		param PostProjectLikeParam
		err   error
	)

	if err = gin2.ShouldBind(c, &param); err != nil {
		common.ReturnError(c, &dberror.Error{Code: ecode.RequestErr, Message: err.Error()})
		return
	}
	account := c.GetString("account")
	if account == "" {
		common.ReturnError(c, &dberror.Error{Code: ecode.RequestErr, Message: "account error"})
		return
	}

	err = service.API.PostProjectLike(account, param.ID)
	if err != nil {
		common.ReturnError(c, err)
		return
	}
	c.JSON(ecode.OK, common.Resp(ecode.OK, nil))
}

type PostProjectUnLikeParam struct {
	ID uint64 `json:"id" form:"id" binding:"required"`
}

// postProjectUnLike un like项目
// @Summary un like项目
// @Tags 项目
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param object query PostProjectLikeParam false "查询参数"
// @Router /project/un/like [post]
func postProjectUnLike(c *gin.Context) {
	var (
		param PostProjectUnLikeParam
		err   error
	)

	if err = gin2.ShouldBind(c, &param); err != nil {
		common.ReturnError(c, &dberror.Error{Code: ecode.RequestErr, Message: err.Error()})
		return
	}
	account := c.GetString("account")
	if account == "" {
		common.ReturnError(c, &dberror.Error{Code: ecode.RequestErr, Message: "account error"})
		return
	}

	err = service.API.PostProjectUnLike(account, param.ID)
	if err != nil {
		common.ReturnError(c, err)
		return
	}
	c.JSON(ecode.OK, common.Resp(ecode.OK, nil))
}

type PostProjectSuperLikeParam struct {
	ID uint64 `json:"id" form:"id" binding:"required"`
}

// postProjectSuperLike super like项目
// @Summary super like项目
// @Tags 项目
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param object query PostProjectSuperLikeParam false "查询参数"
// @Router /project/super/like [post]
func postProjectSuperLike(c *gin.Context) {
	var (
		param PostProjectSuperLikeParam
		err   error
	)

	if err = gin2.ShouldBind(c, &param); err != nil {
		common.ReturnError(c, &dberror.Error{Code: ecode.RequestErr, Message: err.Error()})
		return
	}
	account := c.GetString("account")
	if account == "" {
		common.ReturnError(c, &dberror.Error{Code: ecode.RequestErr, Message: "account error"})
		return
	}

	err = service.API.PostProjectSuperLike(account, param.ID)
	if err != nil {
		common.ReturnError(c, err)
		return
	}
	c.JSON(ecode.OK, common.Resp(ecode.OK, nil))
}
