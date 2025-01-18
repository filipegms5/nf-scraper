package main

import (
	"os"

	"github.com/filipegms5/nf-scraper/router"
)

func main() {
	router := router.SetupRouter()
	router.Run(os.Getenv("PORT"))
}
