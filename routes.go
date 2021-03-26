package main

import (
	"github.com/gin-gonic/gin"
	"oceanlearn.teach/ginessential/controller"
	"oceanlearn.teach/ginessential/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	// r.GET("/api/auth/info", controller.Info) // 不使用中间件保护的情况，返回的数据中user为null，因为上下文没有存user的信息
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	return r
}
