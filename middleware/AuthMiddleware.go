package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jkdev.cn/api/common"
	"jkdev.cn/api/model"
	"net/http"
	"strings"
)

//这些中间件，相当于 java 里的 filter，这里用来进行权限认证
/*
 filter 的执行逻辑为

	httpContext(包括 request 和 response)  ->

									filter1  ->  filter2  ->  ...
																   ->  controller
																				   ->  ...  -> filer2  -> filer1

																											->   httpContext(包括 request 和 response)
一条 http 消息，会经过 过滤器链， controller，在过滤器链的 末端， controller 处理之后， 返回时 会倒序 经过 过滤器链

过滤器链可以 解耦某些功能，比如 鉴权、黑名单、打印请求日志、获取请求执行时间
可以不改变原来的 controller 的逻辑的基础上，添加新的功能

鉴权时，如果判断不满足条件，可以提前返回，就不用走后面的流程了
*/

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取 authorization header
		tokenString := ctx.GetHeader("Authorization")

		fmt.Print("请求token", tokenString)

		//validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:] //截取字符

		token, claims, err := common.ParseToken(tokenString)

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		//token通过验证, 获取claims中的UserID
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		// 验证用户是否存在
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		//用户存在 将user信息写入上下文
		ctx.Set("user", user)

		// 这里就是 执行下一个 filter，filter 的尽头就是 controller
		ctx.Next()

		//这里写代码的话，就是在 controller 处理完之后才会走到这里了
	}
}
