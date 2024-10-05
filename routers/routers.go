package routers

import (
	"clock/controller"
	"clock/controller/userHandler"
	"clock/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	r.POST("/user/register", userHandler.Register)
	r.POST("/user/login", userHandler.Login)

	public := r.Group("/public")
	public.Use(middleware.AuthMiddleware())
	public.POST("/start", controller.Start)
	public.POST("/end", controller.End)
	public.GET("/rank", controller.Rank)
	public.GET("/userInfo", userHandler.UserInfo)
	return r
}
