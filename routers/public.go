package routers

import (
	"clock/controller"
	"clock/controller/userHandler"
	"clock/middleware"

	"github.com/gin-gonic/gin"
)

func RouterGroup_User(r *gin.Engine) {
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware())
	user.POST("/start", controller.Start)
	user.POST("/end", controller.End)
	user.GET("/rank", controller.Rank)
	user.GET("/userInfo", userHandler.UserInfo)
}
