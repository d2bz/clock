package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		cookie, err := context.Request.Cookie("Telephone")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				context.JSON(http.StatusUnauthorized, gin.H{
					"code": 401,
					"msg":  "未找到用户信息",
				})
				context.Abort()
				return
			}
			context.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "服务器错误",
			})
			context.Abort()
			return
		}
		userTel := cookie.Value
		context.Set("user_tel", userTel)
		context.Next()
	}
}
