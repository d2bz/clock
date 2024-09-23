package routes

import (
	"clock/controller"
	"clock/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	r.POST("/user/register", controller.Register)
	r.POST("/user/login", controller.Login)

	public := r.Group("/public")
	public.Use(middleware.AuthMiddleware())
	public.POST("/start", controller.Start)
	public.POST("/end", controller.End)
	public.GET("/rank", controller.Rank)
	return r
}
