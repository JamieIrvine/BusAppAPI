package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type TranslinkLiveBus struct {
	ID                string `json:"ID"`
	Operator          string `json:"Operator"`
	JourneyIdentifier string `json:"JourneyID"`
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

type LiveBusOutput struct {
	ID                   string `json:"id"`
	Operator             string `json:"operator"`
	JourneyID            string `json:"journeyID"`
	DateOfOperation      string `json:"dateOfOperation"`
	DelaySeconds         int    `json:"delaySeconds"`
	DelaySecondsReadable string `json:"delaySecondsReadable"`
	LastUpdated          string `json:"lastUpdated"`
	LastUpdatedPrevious  string `json:"lastUpdatedPrevious"`
	XCurrent             string `json:"xCurrent"`
	YCurrent             string `json:"yCurrent"`
	XPrevious            string `json:"xPrevious"`
	YPrevious            string `json:"yPrevious"`
	VehicleID            string `json:"vehicleID"`
	RealtimeAvailable    int    `json:"realtimeAvailable"`
	ServiceNumber        string `json:"serviceNumber"`
	Direction            string `json:"direction"`
}

func convertLiveBuses(buses []TranslinkLiveBus) []LiveBusOutput {
	result := make([]LiveBusOutput, 0, len(buses))
	for _, b := range buses {
		var dateStr string
		t, err := time.Parse("02.01.2006", b.DayOfOperation)
		if err != nil {
			log.Printf("Error parsing day of operation date: %v", err)
			dateStr = ""
		} else {
			dateStr = t.Format("2006-01-02")
		}

		result = append(result, LiveBusOutput{
			ID:                   b.ID,
			Operator:             b.Operator,
			JourneyID:            b.JourneyIdentifier,
			DateOfOperation:      dateStr,
			DelaySeconds:         b.Delay,
			DelaySecondsReadable: secondsToMinutes(b.Delay),
			XCurrent:             b.X,
			YCurrent:             b.Y,
			LastUpdated:          b.Timestamp,
			XPrevious:            b.XPrevious,
			YPrevious:            b.YPrevious,
			LastUpdatedPrevious:  b.TimestampPrevious,
			VehicleID:            b.VehicleIdentifier,
			RealtimeAvailable:    b.RealtimeAvailable,
			ServiceNumber:        b.LineText,
			Direction:            b.DirectionText,
		})
	}
	return result
}

func filterRecentBuses(buses []LiveBusOutput, window time.Duration) []LiveBusOutput {
	now := time.Now()
	out := make([]LiveBusOutput, 0, len(buses))
	for _, b := range buses {
		// Parse RFC3339 like "2025-08-15T22:47:20+01:00"
		t, err := time.Parse(time.RFC3339, b.LastUpdated)
		if err != nil {
			// Skip entries with bad timestamps, but log for visibility.
			log.Printf("skipping bus %s: bad lastUpdated %q: %v", b.ID, b.LastUpdated, err)
			continue
		}
		age := now.Sub(t)
		if age >= 0 && age <= window {
			out = append(out, b)
		}
	}
	return out
}

func secondsToMinutes(seconds int) string {
	minutes := seconds / 60
	return fmt.Sprintf("%d minutes", minutes)
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

		converted := convertLiveBuses(buses)
		// Keep only buses updated within the last 5 minutes
		//converted = filterRecentBuses(converted, 10*time.Minute)

		serviceNumberParam := c.Query("serviceNumber")
		if serviceNumberParam != "" {
			allBuses := converted
			converted = make([]LiveBusOutput, 0, len(allBuses))
			for _, b := range allBuses {
				if strings.EqualFold(b.ServiceNumber, serviceNumberParam) {
					converted = append(converted, b)
				}
			}
		}

		c.JSON(http.StatusOK, converted)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
