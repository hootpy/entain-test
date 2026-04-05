package db

import (
	"database/sql"
	"errors"
	"strings"
	"sync"
	"time"

	"sports/proto/sports"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EventsRepo provides repository access to events.
type EventsRepo interface {
	// Init will initialise our events repository.
	Init() error

	// List will return a list of events.
	List(filter *sports.ListEventsRequestFilter) ([]*sports.Event, error)
}

type eventsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewEventsRepo creates a new events repository.
func NewEventsRepo(db *sql.DB) EventsRepo {
	return &eventsRepo{db: db}
}

// Init prepares the events repository dummy data.
func (r *eventsRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy events.
		err = r.seed()
	})

	return err
}

func (r *eventsRepo) List(filter *sports.ListEventsRequestFilter) ([]*sports.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getEventQueries()[eventsList]

	query, args = r.applyFilter(query, filter)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanEvents(rows)
}

func (r *eventsRepo) applyFilter(query string, filter *sports.ListEventsRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)
	currentTime := time.Now().UTC()

	args = append(args, currentTime)

	if filter == nil {
		return query, args
	}

	if filter.Visible != nil {
		clauses = append(clauses, "visible = ?")
		args = append(args, *filter.Visible)
	}

	if filter.Status != nil {
		clauses = append(clauses, "status = ?")
		args = append(args, strings.ToUpper(*filter.Status))
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	order := "time_asc"
	if filter.Order != nil {
		order = strings.ToLower(*filter.Order)
	}

	switch order {
	case "name_asc":
		query += " ORDER BY name ASC"
	case "name_desc":
		query += " ORDER BY name DESC"
	case "time_desc":
		query += " ORDER BY advertised_start_time DESC"
	default:
		query += " ORDER BY advertised_start_time ASC"
	}

	return query, args
}

func (r *eventsRepo) scanEvents(
	rows *sql.Rows,
) ([]*sports.Event, error) {
	var events []*sports.Event

	for rows.Next() {
		var event sports.Event
		var advertisedStart time.Time
		var status string

		if err := rows.Scan(&event.Id, &event.Name, &advertisedStart, &event.Visible, &status); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}

			return nil, err
		}

		ts := timestamppb.New(advertisedStart)

		event.AdvertisedStartTime = ts
		event.Status = status

		events = append(events, &event)
	}

	return events, nil
}
