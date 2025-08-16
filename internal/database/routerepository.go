package database

import (
	"bus-app-api/internal/models"
	"database/sql"
)

type RouteRepository struct {
	db *sql.DB
}

func NewRouteRepository(db *sql.DB) RouteRepository {
	return RouteRepository{db: db}
}

func (r *RouteRepository) GetAll() ([]models.Route, error) {

	var routes []models.Route
	query := `
        SELECT id, agency_id, service_number, route_name, route_type, direction FROM routes
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var route models.Route
		err := rows.Scan(&route.ID, &route.AgencyID, &route.ServiceNumber, &route.RouteName, &route.RouteType, &route.Direction)
		if err != nil {
			return nil, err
		}

		routes = append(routes, route)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

func (r *RouteRepository) Upsert(routes []models.Route) error {

	transaction, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO routes (id, agency_id, service_number, route_name, route_type, direction) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (service_number, route_name) DO UPDATE SET
		    agency_id = EXCLUDED.agency_id,
			service_number = EXCLUDED.service_number,
			route_name = EXCLUDED.route_name,
			route_type = EXCLUDED.route_type,
			direction = EXCLUDED.direction
	`

	for _, route := range routes {
		_, err := ExecTransaction(transaction, query, route.ID, route.AgencyID, route.ServiceNumber, route.RouteName, route.RouteType, route.Direction)
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
