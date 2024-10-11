package middleware

import (
	"clock/common"
	"clock/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// jwt中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		//获取authorization header
		tokenString := context.Request.Header.Get("Authorization")
		//token为空
		if tokenString == "" {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token为空",
			})
			context.Abort()
			return
		}
		//错误token
		if len(tokenString) < 100 || !strings.HasPrefix(tokenString, "Bearer") {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "错误token",
			})
			context.Abort()
			return
		}
		//提取token的有效部分，去掉Bearer前缀
		tokenString = tokenString[7:]
		//解析token
		token, claims, err := common.ParseToken(tokenString)
		//非法token
		if err != nil || !token.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "非法token",
			})
			context.Abort()
			return
		}
		//获取claims中的信息
		uid := claims.UID
		db := common.DB
		var curUser model.User
		db.Where("user_id = ?", uid).First(&curUser)
		//写入上下文
		context.Set("curUser", curUser)
		context.Next()
	}
}

//cookie中间件
/*
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
*/
