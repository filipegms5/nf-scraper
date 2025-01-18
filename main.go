package main

import (
	"os"

	"github.com/filipegms5/nf-scraper/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000" // Default port if not specified
	}
	router := router.SetupRouter()
	err := router.Run(port)
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
		return
	}
}
