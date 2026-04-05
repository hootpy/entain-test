package db

import (
	"database/sql"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error)

	// Get will return a race by ID.
	Get(req *racing.GetRaceRequest) (*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) Get(req *racing.GetRaceRequest) (*racing.Race, error) {
	query := getRaceQueries()[raceGet]
	row := r.db.QueryRow(query, time.Now().UTC(), req.Id)

	return r.scanRace(row)
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)
	currentTime := time.Now().UTC()

	args = append(args, currentTime)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
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
	case "meeting_asc":
		query += " ORDER BY meeting_id ASC"
	case "meeting_desc":
		query += " ORDER BY meeting_id DESC"
	case "time_desc":
		query += " ORDER BY advertised_start_time DESC"
	default:
		query += " ORDER BY advertised_start_time ASC"
	}

	return query, args
}

func (r *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time
		var status string

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart, &status); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts
		race.Status = status

		races = append(races, &race)
	}

	return races, nil
}

func (r *racesRepo) scanRace(row *sql.Row) (*racing.Race, error) {
	race := &racing.Race{}
	var advertisedStart time.Time
	var status string

	if err := row.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart, &status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	ts, err := ptypes.TimestampProto(advertisedStart)
	if err != nil {
		return nil, err
	}

	race.AdvertisedStartTime = ts
	race.Status = status

	return race, nil
}
