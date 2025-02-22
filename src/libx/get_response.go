package libx

import (
	"github.com/gin-gonic/gin"
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
func Ok(c *gin.Context, input ...interface{}) {
	if len(input) >= 3 {
		log.Println("too many parameters")
		Err(c, 500, "参数过多，请后端开发人员排查", nil)
	}
	Code(c, 200)
	if len(input) == 2 {
		Msg(c, input[0].(string))
		Data(c, input[1])
	} else {
		Msg(c, input[0].(string))
		Data(c, nil)
	}
}

// 一个参数省略msg
func Registered(c *gin.Context, input ...interface{}) {
	if len(input) >= 3 {
		log.Println("too many parameters")
		Err(c, 500, "参数过多，请后端开发人员排查", nil)
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

// 一个参数省略msg。
func Fail(c *gin.Context, input ...interface{}) {
	if len(input) >= 3 {
		log.Println("too many parameters")
		Err(c, 500, "参数过多，请后端开发人员排查", nil)
	}
	Code(c, 400)
	if len(input) == 2 {
		Msg(c, input[0].(string))
		Data(c, input[1])
	} else {
		Msg(c, "fail")
		Data(c, input[0])
	}
}

func Err(c *gin.Context, code int, msg string, err error) {
	//Code(c, code)
	// 如果4开头的错误码（6位），就返回400
	// 4开头代表业务错误，5开头代表我的错误
	if code >= 400000 && code < 500000 {
		Code(c, 400)
	} else {
		Code(c, 500)
	}

	var errorMsg string
	if err != nil {
		errorMsg = err.Error()
	}
	c.Set("code", code)
	c.Set("message", msg+" "+errorMsg)
	// 打印错误信息
	log.Println(msg + " " + errorMsg)

	//Data(c, gin.H{
	//	"error": errorMsg,
	//})
}
