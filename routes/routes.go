package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kewlspy/web-backend/controllers"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/urls", controllers.SubmitURL)
		api.GET("/urls", controllers.GetAllURLs)
		api.GET("/urls/:id", controllers.GetURLDetails)
	}
}
