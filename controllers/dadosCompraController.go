package controllers

import (
	"net/http"

	"github.com/filipegms5/nf-scraper/services"

	"github.com/gin-gonic/gin"
)

type RequestPayload struct {
	URL string `json:"url" binding:"required"`
}

func FetchDadosCompra(c *gin.Context) {
	var payload RequestPayload

	// Bind the JSON payload to the struct
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the URL variable
	url := payload.URL

	dadosCompra, err := services.FetchDadosCompra(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dadosCompra)
}
