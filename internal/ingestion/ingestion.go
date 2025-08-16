package ingestion

import (
	"bus-app-api/internal/database"
	"bus-app-api/internal/models"
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func Ingest(db *sql.DB) error {
	err := ingestBusStops(db)
	if err != nil {
		return err
	}

	log.Println("Completed ingestion of bus stop data")
	return nil
}

func ingestBusStops(db *sql.DB) error {
	busStopsFilePath := "./ingestion/Busstops17-06--2025.csv"

	file, err := os.Open(busStopsFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip the header row first
	_, err = reader.Read()
	if err != nil {
		return err
	}

	stops := make([]models.Stop, 0)
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		latitude, err := strconv.ParseFloat(record[8], 64)
		if err != nil {
			return err
		}

		longitude, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			return err
		}

		stop := models.Stop{
			ID:        record[0],
			Name:      record[2],
			Latitude:  latitude,
			Longitude: longitude,
		}

		stops = append(stops, stop)
	}

	stopRepository := database.NewStopRepository(db)
	return stopRepository.Upsert(stops)
}
