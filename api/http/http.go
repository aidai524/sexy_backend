package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"net/http"
	"sexy_backend/api/conf"
	g "sexy_backend/api/gin"
	"sexy_backend/common/ecode"
	common "sexy_backend/common/http"
	_ "sexy_backend/docs" // 导入生成的文档包
	"strconv"
	"time"
)

func Init() {
	initRouter()
}

func initRouter() {
	r := gin.New()
	r.Use(g.Recovery())
	r.Use(g.Log())
	if conf.Conf.Cors {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     conf.Conf.AllowOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "locale", "Authorization"},
			AllowCredentials: true,
			MaxAge:           24 * time.Hour,
		}))
	}
	if conf.Conf.Debug {
		r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	}

	app := r.Group("/api/v1")
	{
		app.GET("/status", ping)
		app.GET("/project", getProject)              // 获取项目
		app.GET("/project/search", getProjectSearch) // 搜索项目
		app.GET("/project/list", getProjectList)     // 获取项目列表
	}

	token := r.Group("/api/v1")
	{
		token.GET("/account/token", getAccountToken) // 获取请求token
	}

	auth := r.Group("/api/v1")
	auth.Use(Auth())
	{
		auth.POST("/project", postProject)                     // 上传项目
		auth.POST("/project/like", postProjectLike)            // like项目
		auth.POST("/project/un/like", postProjectUnLike)       // un like项目
		auth.POST("/project/super/like", postProjectSuperLike) // super like项目
	}

	err := r.Run(":" + strconv.Itoa(conf.Conf.Port))
	if err != nil {
		panic(err)
	}
}

// ping 服务器状态检查
// @Summary ping
// @Tags 系统
// @Accept application/json
// @Produce application/json
// @Success 200 {object} map[string]int64{} "响应结果"
// @Router /status [get]
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, common.Resp(ecode.OK, map[string]interface{}{
		"timestamp": time.Now().Unix(),
	}))
}
