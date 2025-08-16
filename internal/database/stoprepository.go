package database

import (
	"bus-app-api/internal/models"
	"database/sql"
)

type StopRepository struct {
	db *sql.DB
}

func NewStopRepository(db *sql.DB) StopRepository {
	return StopRepository{db: db}
}

func (r *StopRepository) GetAll() ([]models.Stop, error) {

	var stops []models.Stop
	query := `
        SELECT id, name, latitude, longitude FROM stops
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var stop models.Stop
		err := rows.Scan(&stop.ID, &stop.Name, &stop.Latitude, &stop.Longitude)
		if err != nil {
			return nil, err
		}

		stops = append(stops, stop)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stops, nil
}

// Upsert inserts or updates stops based on their ID
func (r *StopRepository) Upsert(stops []models.Stop) error {

	transaction, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Upsert query: insert new stops or update existing ones based on ID
	query := `
		INSERT INTO stops (id, name, latitude, longitude) VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude
	`

	for _, stop := range stops {
		_, err := ExecTransaction(transaction, query, stop.ID, stop.Name, stop.Latitude, stop.Longitude)
		if err != nil {
			return err
		}
	}

	err = transaction.Commit()
	if err != nil {
		return err
	}

	return nil
}
