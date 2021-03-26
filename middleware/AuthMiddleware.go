package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"oceanlearn.teach/ginessential/common"
	"oceanlearn.teach/ginessential/model"
)

// gin的中间件就是一个函数，然后返回一个HandlerFunc
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取authorization header
		tokenString := ctx.GetHeader("Authorization")
		// validate token format
		// 如果没传token，或者token格式错误就返回权限不足
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			// 丢弃本次请求并返回
			ctx.Abort()
			return
		}
		// 如果token验证有效，则提取token的有效部分
		// "Bearer "占了7位
		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid { // 如果解析失败，或者token无效
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}
		// 验证通过后获取claims中的user.Id
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)
		// 验证用户是否存在，如果用户不存在，则token无效
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}
		// 如果用户存在，将user信息写入上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}
