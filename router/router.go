package router

import (
	"github.com/filipegms5/nf-scraper/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/dados-compra", func(c *gin.Context) {
		controllers.FetchDadosCompra(c)
	})

	return router
}
