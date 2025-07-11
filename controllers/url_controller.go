package controllers

import (
	"net/http"
	"github.com/kewlspy/web-backend/models"
	"github.com/kewlspy/web-backend/services"

	"github.com/gin-gonic/gin"
)

func SubmitURL(c *gin.Context) {
	var input struct {
		Link string `json:"url" binding:"required,url"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url := models.URL{Link: input.Link, Status: "queued"}
	models.DB.Create(&url)

	go services.CrawlURL(&url)

	c.JSON(http.StatusAccepted, gin.H{"message": "Crawl started", "id": url.ID})
}

func GetAllURLs(c *gin.Context) {
	var urls []models.URL
	models.DB.Find(&urls)
	c.JSON(http.StatusOK, urls)
}

func GetURLDetails(c *gin.Context) {
	var url models.URL
	if err := models.DB.Preload("BrokenLinkItems").First(&url, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.JSON(http.StatusOK, url)
}
