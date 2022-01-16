package events

import (
	"context"
	"errors"
	"pkg/pkg/models"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type DBManager interface {
	EventExists(event models.Event, tx pgx.Tx) (bool, error)
	InsertEvent(event models.Event, tx pgx.Tx) error
}

type postgresDBManager struct {
	db *pgx.Conn
}

func NewEventsDBManager(db *pgx.Conn) DBManager {
	return postgresDBManager{
		db: db,
	}
}

func (p postgresDBManager) EventExists(event models.Event, tx pgx.Tx) (bool, error) {
	var id string
	if err := tx.QueryRow(context.Background(), "select id from events.processed_events where id=$1 and event_name=$2", event.ID(), event.Name()).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.WithField("code", pgErr.Code).
				WithField("message", pgErr.Message).
				Error("encountered an issue querying for the event")
			return false, err
		}
		return false, nil
	}

	return len(id) > 0, nil
}

// InsertEvent method for create a new Event.
// @Description Insert a new Event in DB.
// @Summary insert a new Event in DB.
// @Tags event DB
// @Accept json
// @Produce json
// @Success 200 {string} status "ok"
// @Router /v1/token/new [get]
func (p postgresDBManager) InsertEvent(event models.Event, tx pgx.Tx) error {
	if _, err := tx.Exec(context.Background(), "insert into events.processed_events (id, event_name, processed_timestamp) values ($1, $2, $3)", event.ID(), event.Name(), time.Now()); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.WithField("code", pgErr.Code).
				WithField("message", pgErr.Message).
				Error("encountered an issue inserting the event into the DB")
		}

		return err
	}
	return nil
}
