package main

import (
	"github.com/kewlspy/web-backend/models"
	"github.com/kewlspy/web-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDB() // ✅ Don't forget this line!
	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run(":8080") // run on localhost:8080
}
