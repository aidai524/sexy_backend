package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	g "github.com/gin-gonic/gin"
	"go.uber.org/atomic"
	"reflect"
	"sexy_backend/common/ecode"
	"sexy_backend/common/log"
	model2 "sexy_backend/common/model"
	"sexy_backend/common/sexyerror"
	"strings"
	"time"
)

var findClientIPError = atomic.NewUint64(0)

type Param struct {
	Limit int64 `json:"limit"`
}

func ShouldBind(c *g.Context, obj interface{}) (err error) {
	return ShouldBindMaxLimit(c, obj, 100)
}

func ShouldBindMaxLimit(c *g.Context, obj interface{}, limitMax int64) (err error) {
	var max = limitMax
	err = c.ShouldBind(obj)
	if err != nil {
		return
	}

	var (
		v = reflect.ValueOf(obj).Elem()
	)
	for v.Kind() == reflect.Slice || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		log.Error("from and to must be pointer")
		err = fmt.Errorf("from and to must be pointer")
		return
	}

	if !v.IsValid() {
		err = &sexyerror.Error{Code: ecode.UnknownError, Message: fmt.Sprintf("param error")}
		return
	}

	if obj != nil {
		p := &Param{}
		err = model2.Copy(p, obj)
		if err != nil {
			return
		}
		if p.Limit > max {
			err = &sexyerror.Error{Code: ecode.ParamLimitMaxError, Message: fmt.Sprintf("param limit max %v", max)}
			return
		}
	}
	err = model2.Validate(obj)
	if err != nil {
		return
	}
	return
}

func MaxStartTime(start int64, findStartTime int) int64 {
	t := time.Now().UTC().AddDate(0, 0, -findStartTime)
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	if start < 9999999999 {
		// 秒
		if t.Unix() < start {
			return start
		} else {
			return t.Unix()
		}
	} else {
		// 毫秒
		if t.UnixNano()/1e6 < start {
			return start
		} else {
			return t.UnixNano() / 1e6
		}
	}
}

func GetRealIP(c *gin.Context) string {
	IPValue := c.Request.Header.Get("X-Original-Forwarded-For")
	IPList := strings.Split(IPValue, ",")
	if l := len(IPList); l > 0 {
		return strings.TrimSpace(IPList[l-1])
	}

	findClientIPError.Add((1))
	log.Error("No sufficient IPs found in Header for key X-Original-Forwarded-For,  IPValue : %v", IPValue)
	return ""
}

func GetFindClientIPErrorCount() uint64 {
	return findClientIPError.Load()
}
