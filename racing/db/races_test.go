package db

import (
	"database/sql"
	"testing"
	"time"

	"git.neds.sh/matty/entain/racing/proto/racing"
	_ "github.com/mattn/go-sqlite3"
)

func TestRacesRepo_List_NoFilter(t *testing.T) {
	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer func() { _ = dbConn.Close() }()

	_, err = dbConn.Exec(`
		CREATE TABLE races (
			id INTEGER PRIMARY KEY,
			meeting_id INTEGER,
			name TEXT,
			number INTEGER,
			visible INTEGER,
			advertised_start_time DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("create table error = %v", err)
	}

	start1 := time.Date(2026, 3, 31, 10, 0, 0, 0, time.UTC)
	start2 := time.Date(2026, 3, 31, 11, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO races (id, meeting_id, name, number, visible, advertised_start_time)
		VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)
	`,
		1, 100, "Race One", 1, true, start1,
		2, 200, "Race Two", 2, false, start2,
	)
	if err != nil {
		t.Fatalf("insert races error = %v", err)
	}

	repo := NewRacesRepo(dbConn)

	got, err := repo.List(nil)
	if err != nil {
		t.Fatalf("List(nil) error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("List(nil) returned %d races, want 2", len(got))
	}

	if got[0].Id != 1 {
		t.Fatalf("got[0].Id = %d, want 1", got[0].Id)
	}

	if got[0].MeetingId != 100 {
		t.Fatalf("got[0].MeetingId = %d, want 100", got[0].MeetingId)
	}

	if got[0].Name != "Race One" {
		t.Fatalf("got[0].Name = %q, want %q", got[0].Name, "Race One")
	}

	if got[0].AdvertisedStartTime == nil {
		t.Fatal("got[0].AdvertisedStartTime is nil")
	}

	if got[1].Id != 2 {
		t.Fatalf("got[1].Id = %d, want 2", got[1].Id)
	}
}

func TestRacesRepo_List_Status(t *testing.T) {
	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer func() { _ = dbConn.Close() }()

	_, err = dbConn.Exec(`
		CREATE TABLE races (
			id INTEGER PRIMARY KEY,
			meeting_id INTEGER,
			name TEXT,
			number INTEGER,
			visible INTEGER,
			advertised_start_time DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("create table error = %v", err)
	}

	pastStart := time.Now().UTC().Add(-1 * time.Hour)
	futureStart := time.Now().UTC().Add(1 * time.Hour)

	_, err = dbConn.Exec(`
		INSERT INTO races (id, meeting_id, name, number, visible, advertised_start_time)
		VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)
	`,
		1, 100, "Past Race", 1, true, pastStart,
		2, 200, "Future Race", 2, true, futureStart,
	)
	if err != nil {
		t.Fatalf("insert races error = %v", err)
	}

	repo := NewRacesRepo(dbConn)

	got, err := repo.List(&racing.ListRacesRequestFilter{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("List() returned %d races, want 2", len(got))
	}

	if got[0].Status != "CLOSED" {
		t.Fatalf("got[0].Status = %q, want %q", got[0].Status, "CLOSED")
	}

	if got[1].Status != "OPEN" {
		t.Fatalf("got[1].Status = %q, want %q", got[1].Status, "OPEN")
	}
}

func TestRacesRepo_List_StatusFilter(t *testing.T) {
	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer func() { _ = dbConn.Close() }()

	_, err = dbConn.Exec(`
		CREATE TABLE races (
			id INTEGER PRIMARY KEY,
			meeting_id INTEGER,
			name TEXT,
			number INTEGER,
			visible INTEGER,
			advertised_start_time DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("create table error = %v", err)
	}

	pastStart := time.Now().UTC().Add(-1 * time.Hour)
	futureStart := time.Now().UTC().Add(1 * time.Hour)

	_, err = dbConn.Exec(`
		INSERT INTO races (id, meeting_id, name, number, visible, advertised_start_time)
		VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)
	`,
		1, 100, "Past Race", 1, true, pastStart,
		2, 200, "Future Race", 2, true, futureStart,
	)
	if err != nil {
		t.Fatalf("insert races error = %v", err)
	}

	status := "closed"
	repo := NewRacesRepo(dbConn)

	got, err := repo.List(&racing.ListRacesRequestFilter{Status: &status})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("List(status=closed) returned %d races, want 1", len(got))
	}

	if got[0].Id != 1 {
		t.Fatalf("got[0].Id = %d, want 1", got[0].Id)
	}

	if got[0].Status != "CLOSED" {
		t.Fatalf("got[0].Status = %q, want %q", got[0].Status, "CLOSED")
	}
}
