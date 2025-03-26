package libx

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/trancecho/open-sdk/logx"
	"go.uber.org/zap"
	"log"
)

func Code(c *gin.Context, code int) {
	// 设置 HTTP 状态码
	c.Status(code)
}

func Msg(c *gin.Context, msg string) {
	c.Set("message", msg)
}

func Data(c *gin.Context, data interface{}) {
	c.Set("data", data)
}

// 一个参数省略msg
func Ok(c *gin.Context, msg string, data interface{}) {
	Code(c, 200)  // 设置响应状态码
	Msg(c, msg)   // 设置响应消息
	Data(c, data) // 设置响应数据

	// 日志事件
	logEvent := c.GetString("log_event")
	if logEvent != "" {
		var logData string

		// 检查 data 的类型
		switch d := data.(type) {
		case string:
			logData = d
		case map[string]interface{}: // 处理 gin.H
			bytes, err := json.Marshal(d) // 将 map 序列化为 JSON 字符串
			if err != nil {
				logData = "无法序列化 data"
			} else {
				logData = string(bytes)
			}
		default:
			logData = fmt.Sprintf("%v", d) // 对其他类型调用默认格式化
		}

		logx.Info(logEvent,
			zap.String("data", logData), // 使用安全的日志字符串
		)
	}
}

// 一个参数省略msg
func Registered(c *gin.Context, input ...interface{}) {
	if len(input) >= 3 {
		log.Println("too many parameters")
		Err(c, 500, "参数过多，请后端开发人员排查", ErrOptions{})
	}
	Code(c, 201)
	if len(input) == 2 {
		Msg(c, input[0].(string))
		Data(c, input[1])
	} else {
		Msg(c, input[0].(string))
		Data(c, nil)
	}
}

type ErrOptions struct {
	Err  error `json:"err,omitempty"`
	Code int   `json:"code,omitempty"`
}

func Err(c *gin.Context, code int, msg string, options ErrOptions) {
	Msg(c, msg)
	Code(c, code)
	if options != (ErrOptions{}) {
		Data(c, options)
	}
	if options.Err != nil {
		log.Println("error:", options.Err)
	}
}
