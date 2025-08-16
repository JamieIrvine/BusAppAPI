package database

import (
	"bus-app-api/internal/config"
	"database/sql"
	"fmt"
	"log"
	"net/url"

	_ "github.com/lib/pq"
)

func buildConnectionString(cfg config.Config) string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s%s",
		cfg.DbConfig.Driver, cfg.DbConfig.Username, url.QueryEscape(cfg.DbConfig.Password), cfg.DbConfig.Host, cfg.DbConfig.Port, cfg.DbConfig.DatabaseName, cfg.DbConfig.ExtraParams)
}

func GetDb(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", buildConnectionString(cfg))
	if err != nil {
		log.Printf("failed to open database: %v", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Printf("database is unreachable: %v", err)
		return nil, err
	}

	return db, nil
}

func RunQuery(cfg config.Config, query string, params ...any) (*sql.Rows, error) {
	db, err := GetDb(cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query, params...)

	if err != nil {
		return nil, err
	}
	return rows, nil
}

func ExecTransaction(transaction *sql.Tx, query string, params ...any) (*sql.Tx, error) {
	var err error
	if len(params) > 0 {
		_, err = transaction.Exec(query, params...)
	} else {
		_, err = transaction.Exec(query)
	}

	if err != nil {
		transaction.Rollback()
		return transaction, err
	}

	return transaction, nil
}
