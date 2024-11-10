package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"sexy_backend/common/ecode"
	commonHttp "sexy_backend/common/http"
	"sexy_backend/common/stack"
	"strings"
)

func Recovery() func(*gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					c.Abort()
				} else {
					_ = c.Error(fmt.Errorf("%v, stack: %v", err, stack.Stack(2, 1000)))
					c.AbortWithStatusJSON(http.StatusInternalServerError,
						commonHttp.RespWithMsg(ecode.ServerErr, "internal server error"))
					writeLog(c)
				}
			}
		}()
		c.Next()
	}
}
