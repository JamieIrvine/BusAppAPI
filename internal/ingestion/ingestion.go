package ingestion

import (
	"bus-app-api/internal/database"
	"bus-app-api/internal/models"
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
)

func Ingest(db *sql.DB) error {
	err := ingestBusStops(db)
	if err != nil {
		return err
	}
	log.Println("Completed ingestion of bus stop data")

	err = ingestRoutes(db)
	if err != nil {
		return err
	}
	log.Println("Completed ingestion of route data")
	return nil
}

func ingestBusStops(db *sql.DB) error {
	filePath := "./ingestion/Busstops17-06--2025.csv"

	file, err := os.Open(filePath)
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

func ingestRoutes(db *sql.DB) error {
	filePath := "./ingestion/gtfs/metro/routes.txt"

	file, err := os.Open(filePath)
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

	routes := make([]models.Route, 0)
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		routeType := "Unknown"
		if record[4] == "3" {
			routeType = "Bus"
		}

		routeNameSplit := strings.Split(record[3], " | ")
		for i, routeName := range routeNameSplit {
			route := models.Route{
				ID:            record[0],
				AgencyID:      record[1],
				ServiceNumber: record[2],
				RouteName:     routeName,
				RouteType:     routeType,
			}

			if len(routeNameSplit) == 2 {
				// There are two directions. Outbound should be first.
				direction := "Outbound"
				if i == 1 {
					direction = "Inbound"
				}
				route.Direction = direction
			}

			routes = append(routes, route)
		}
	}

	routeRepository := database.NewRouteRepository(db)
	return routeRepository.Upsert(routes)
}
