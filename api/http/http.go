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
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "locale"},
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
