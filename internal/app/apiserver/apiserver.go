package apiserver

import (
	"database/sql"
	_ "github.com/lib/pq"
	"net/http"
	"rest_api/internal/app/store/sqlstore"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	srv := newServer(store)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(sqlString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", sqlString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
