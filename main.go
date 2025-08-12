package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TranslinkLiveBus struct {
	ID                string `json:"id"`
	Operator          string `json:"Operator"`
	JourneyIdentifier string `json:"JourneyIdentifier"`
	DayOfOperation    string `json:"DayOfOperation"`
	Delay             int    `json:"Delay"`
	MOTCode           int    `json:"MOTCode"`
	X                 string `json:"X"`
	Y                 string `json:"Y"`
	Timestamp         string `json:"Timestamp"`
	XPrevious         string `json:"XPrevious"`
	YPrevious         string `json:"YPrevious"`
	TimestampPrevious string `json:"TimestampPrevious"`
	VehicleIdentifier string `json:"VehicleIdentifier"`
	RealtimeAvailable int    `json:"RealtimeAvailable"`
	LineText          string `json:"LineText"`
	DirectionText     string `json:"DirectionText"`
}

func main() {

	r := gin.Default()

	r.GET("/live-buses", func(c *gin.Context) {
		translinkLiveAPIURL := "https://vpos.translinkniplanner.co.uk/velocmap/vmi/VMI"

		resp, err := http.Get(translinkLiveAPIURL)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer resp.Body.Close()

		var buses []TranslinkLiveBus
		if err := json.NewDecoder(resp.Body).Decode(&buses); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, buses)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
