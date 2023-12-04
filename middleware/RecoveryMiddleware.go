package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jkdev.cn/api/response"
)

//返回 "请求异常" 的响应

func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.Fail(ctx, fmt.Sprint(err), nil)
			}
		}()

		ctx.Next()
	}
}
