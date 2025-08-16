package main

import (
	"bus-app-api/internal/config"
	"bus-app-api/internal/database"
	"bus-app-api/internal/endpoints/live"
	"bus-app-api/internal/ingestion"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	db, err := database.GetDb(cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = ingestion.Ingest(db)
	if err != nil {
		log.Fatalf(err.Error())
	}

	r := gin.Default()

	r.GET("/live-buses", func(c *gin.Context) {
		live.Get(c)
	})

	r.GET("/stops", func(c *gin.Context) {
		repo := database.NewStopRepository(db)
		stops, err := repo.GetAll()
		if err != nil {
			log.Fatalf(err.Error())
		}
		c.JSON(http.StatusOK, stops)
	})

	r.GET("/routes", func(c *gin.Context) {
		repo := database.NewRouteRepository(db)
		stops, err := repo.GetAll()
		if err != nil {
			log.Fatalf(err.Error())
		}
		c.JSON(http.StatusOK, stops)
	})

	if err := r.Run(":2877"); err != nil {
		log.Fatal(err)
	}
}
