package controllers

import (
	"net/http"

	"github.com/filipegms5/nf-scraper/services"

	"github.com/gin-gonic/gin"
)

func FetchDadosCompra(c *gin.Context) {
	url := c.Query("url")
	dadosCompra, err := services.FetchDadosCompra(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	c.JSON(http.StatusOK, dadosCompra)
}
