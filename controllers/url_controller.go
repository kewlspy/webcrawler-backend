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

func RetryCrawl(c *gin.Context) {
	id := c.Param("id")
	var url models.URL
	if err := models.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	url.Status = "queued"
	models.DB.Save(&url)

	// Start crawl again asynchronously
	go services.CrawlURL(&url)

	c.JSON(http.StatusOK, url)
}
func DeleteURL(c *gin.Context) {
	id := c.Param("id")

	models.DB.Where("url_id = ?", id).Delete(&models.BrokenLink{})
	models.DB.Delete(&models.URL{}, id)

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
