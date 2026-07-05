package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *Handler) {

	v1 := router.Group("/v1")
	{
		v1.POST("/check", handler.CheckRateLimit)
	}

	admin := v1.Group("/admin")
	{
		admin.POST("/client", handler.CreateClient)
		admin.GET("/client/:clientKey", handler.GetClient)
		admin.PUT("/client/:clientKey", handler.UpdateClient)
		admin.DELETE("/client/:clientKey", handler.DeleteClient)
	}
}