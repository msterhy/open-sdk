package util

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func Ok(c *gin.Context, message string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(200, gin.H{
		"err_code": 200000,
		"message":  message,
		"data":     data,
	})
}

// 规范：错误码提供三位，从000开始
func ClientError(c *gin.Context, errCode int, message string) {
	log.Println("客户端错误", errCode, message)
	c.JSON(400, gin.H{
		"err_code": 400000 + errCode,
		"message":  message,
	})
}

func ServerError(c *gin.Context, errCode int, message string) {
	log.Println("服务端错误", errCode, message)
	c.JSON(500, gin.H{
		"err_code": 500000 + errCode,
		"message":  message,
	})
}

// 处理gRPC错误
func HandleGRPCError(c *gin.Context, err error) {
	st, _ := status.FromError(err)
	c.JSON(grpcToHTTPStatus(st.Code()), gin.H{
		"error":   st.Message(),
		"code":    st.Code().String(),
		"details": st.Details(),
	})
}

// 状态码转换
func grpcToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return 200
	case codes.InvalidArgument:
		return 400
	case codes.NotFound:
		return 404
	// ...其他状态码映射
	default:
		return 500
	}
}
