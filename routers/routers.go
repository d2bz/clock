package routers

import (
	"clock/controller/userHandler"
	"clock/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRouter(r *gin.Engine) {
	r.Use(middleware.CORSMiddleware())
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	RouterGroup_User(r)
}
