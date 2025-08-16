package models

type Route struct {
	ID            string `json:"id"`
	AgencyID      string `json:"agencyId"`
	ServiceNumber string `json:"serviceNumber"`
	RouteName     string `json:"routeName"`
	RouteType     string `json:"routeType"`
	Direction     string `json:"direction"`
}
