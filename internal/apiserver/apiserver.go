package apiserver

import (
	"database/sql"
	"net/http"

	"github.com/ositlar/effective/internal/sqlstore"
	"github.com/sirupsen/logrus"
)

func connectToDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	logrus.Infoln("Connected to database on url ", dbURL)
	return db, nil
}

func StartServer(config *Config) error {
	db, err := connectToDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	store := sqlstore.NewStore(db)
	srv := NewServer(store)
	srv.logger.Infoln("Server started on port", config.BindAddr)
	return http.ListenAndServe(config.BindAddr, srv)
}
