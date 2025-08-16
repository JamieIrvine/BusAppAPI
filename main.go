package main

import (
	"bus-app-api/internal/endpoints/live"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/live-buses", func(c *gin.Context) {
		live.Get(c)
	})

	if err := r.Run(":2877"); err != nil {
		log.Fatal(err)
	}
}
