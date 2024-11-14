package common

import (
	"net/http"
	"sexy_backend/common/ecode"
	"sexy_backend/common/sexyerror"
	"time"

	"github.com/gin-gonic/gin"
)

type PageParam struct {
	PageSize int `form:"page_size" json:"page_size"`
	PageNum  int `form:"page_number" json:"page_number"`
}

var (
	isDev bool
)

func InitHttp(dev bool) {
	isDev = dev
}

func (pr *PageParam) CalcOffsetLimit() (limit, offset int) {
	if pr.PageSize < 1 {
		pr.PageSize = 20
	}
	if pr.PageNum < 1 {
		pr.PageNum = 1
	}
	return pr.PageSize, (pr.PageNum - 1) * pr.PageSize
}

type PageResult struct {
	Data        interface{} `json:"list"`
	Total       int         `json:"total"`
	PageSize    int         `json:"page_size"`
	PageNum     int         `json:"page_number"`
	HasNextPage bool        `json:"has_next_page"`
}

func (p *PageResult) CalcHaxNextPage() {
	p.HasNextPage = p.PageSize*p.PageNum < p.Total
}

func Resp(code int, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code": code,
		"data": data,
	}
}

func RespWithMsg(code int, message string) map[string]interface{} {
	return map[string]interface{}{
		"code":       code,
		"message":    message,
		"time_stamp": time.Now().UnixNano() / 1e6,
	}
}

func GetReturnError(err error) (newError *sexyerror.Error) {
	var statusCode int
	switch e := err.(type) {
	case *sexyerror.Error:
		switch e.Code {
		case ecode.RequestErr:
			statusCode = http.StatusBadRequest
		case ecode.Unauthorized:
			statusCode = http.StatusUnauthorized
		case ecode.Forbidden:
			statusCode = http.StatusForbidden
		default:
			statusCode = http.StatusInternalServerError
		}
		newError = &sexyerror.Error{Code: statusCode, Message: err.Error()}
	case *sexyerror.ThirdPartyError:
		statusCode = http.StatusInternalServerError
		newError = &sexyerror.Error{Code: statusCode, Message: err.Error()}
	default:
		var message string
		if isDev {
			message = err.Error()
		} else {
			message = "Internal Server Error"
		}
		newError = &sexyerror.Error{Code: ecode.ServerErr, Message: message}
	}
	return &sexyerror.Error{Code: ecode.RequestErr, Message: err.Error()}
}

func ReturnError(c *gin.Context, err error) {
	var statusCode int
	switch e := err.(type) {
	case *sexyerror.Error:
		switch e.Code {
		case ecode.RequestErr:
			statusCode = http.StatusBadRequest
		case ecode.Unauthorized:
			statusCode = http.StatusUnauthorized
		case ecode.Forbidden:
			statusCode = http.StatusForbidden
		default:
			statusCode = http.StatusOK
		}
		c.JSON(statusCode, err)
	case *sexyerror.ThirdPartyError:
		statusCode = http.StatusInternalServerError
		c.JSON(statusCode, Resp(ecode.ExternalError, err))
	default:
		statusCode = http.StatusInternalServerError
		var message string
		if isDev {
			message = err.Error()
		} else {
			message = "Internal Server Error"
		}
		c.JSON(statusCode, RespWithMsg(ecode.ServerErr, message))
	}
}

type NewPageResult struct {
	Data  interface{} `json:"list"`
	Total int         `json:"total"`
}

type NewHasNextPageResult struct {
	Data        interface{} `json:"list"`
	HasNextPage bool        `json:"has_next_page"` // 是否有下一页
}
