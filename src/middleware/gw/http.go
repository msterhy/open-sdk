package gw

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine/log"
	"time"
)

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next() // 执行下一个处理程序

		// 打印请求日志
		log.Infof(c, "%s %s %d %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
		)
	}
}

// 自定义请求日志中间件
func requestLogger() {

}
