package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

// DB represents an interaction with a database
type DB struct {
	PostgressDB *pgx.Conn
}

// NewDB returns a default localhost reference to a DB
func NewDB() (DB, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s/%s", "postgress", "postgress", ":5432", "kafka")
	log.WithField("url", url).Info("attempting to connect to DB")

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return DB{}, err
	}

	log.Info("successfully connected to DB")

	return DB{
		PostgressDB: conn,
	}, nil
}
