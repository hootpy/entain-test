package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	_ "github.com/mattn/go-sqlite3"
)

func TestRacingService_ListRaces_VisibleFilter(t *testing.T) {
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
	start3 := time.Date(2026, 3, 31, 12, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO races (id, meeting_id, name, number, visible, advertised_start_time)
		VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)
	`,
		1, 100, "Visible Race 1", 1, true, start1,
		2, 200, "Hidden Race", 2, false, start2,
		3, 300, "Visible Race 2", 3, true, start3,
	)
	if err != nil {
		t.Fatalf("insert races error = %v", err)
	}

	racesRepo := db.NewRacesRepo(dbConn)
	service := NewRacingService(racesRepo)

	// Test filter for visible races only
	visible := true
	request := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{
			Visible: &visible,
		},
	}

	response, err := service.ListRaces(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListRaces() error = %v", err)
	}

	if len(response.Races) != 2 {
		t.Fatalf("ListRaces(visible=true) returned %d races, want 2", len(response.Races))
	}

	// Verify both returned races are visible
	for i, race := range response.Races {
		if !race.Visible {
			t.Fatalf("response.Races[%d].Visible = false, want true", i)
		}
	}

	// Verify we got the correct visible races
	expectedIds := []int64{1, 3}
	for i, race := range response.Races {
		if race.Id != expectedIds[i] {
			t.Fatalf("response.Races[%d].Id = %d, want %d", i, race.Id, expectedIds[i])
		}
	}
}

func TestRacingService_ListRaces_HiddenFilter(t *testing.T) {
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
	start3 := time.Date(2026, 3, 31, 12, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO races (id, meeting_id, name, number, visible, advertised_start_time)
		VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)
	`,
		1, 100, "Visible Race 1", 1, true, start1,
		2, 200, "Hidden Race", 2, false, start2,
		3, 300, "Visible Race 2", 3, true, start3,
	)
	if err != nil {
		t.Fatalf("insert races error = %v", err)
	}

	racesRepo := db.NewRacesRepo(dbConn)
	service := NewRacingService(racesRepo)

	// Test filter for hidden races only
	visible := false
	request := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{
			Visible: &visible,
		},
	}

	response, err := service.ListRaces(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListRaces() error = %v", err)
	}

	if len(response.Races) != 1 {
		t.Fatalf("ListRaces(visible=false) returned %d races, want 1", len(response.Races))
	}

	// Verify the returned race is hidden
	if response.Races[0].Visible {
		t.Fatalf("response.Races[0].Visible = true, want false")
	}

	// Verify we got the correct hidden race
	if response.Races[0].Id != 2 {
		t.Fatalf("response.Races[0].Id = %d, want 2", response.Races[0].Id)
	}
}

func TestRacingService_ListRaces_NoVisibleFilter(t *testing.T) {
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
		1, 100, "Visible Race", 1, true, start1,
		2, 200, "Hidden Race", 2, false, start2,
	)
	if err != nil {
		t.Fatalf("insert races error = %v", err)
	}

	racesRepo := db.NewRacesRepo(dbConn)
	service := NewRacingService(racesRepo)

	// Test without visibility filter (should return all races)
	request := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{
			// No Visible field set - should return all races
		},
	}

	response, err := service.ListRaces(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListRaces() error = %v", err)
	}

	if len(response.Races) != 2 {
		t.Fatalf("ListRaces(no visible filter) returned %d races, want 2", len(response.Races))
	}

	// Verify we got both visible and hidden races
	visibleCount := 0
	hiddenCount := 0
	for _, race := range response.Races {
		if race.Visible {
			visibleCount++
		} else {
			hiddenCount++
		}
	}

	if visibleCount != 1 {
		t.Fatalf("Expected 1 visible race, got %d", visibleCount)
	}

	if hiddenCount != 1 {
		t.Fatalf("Expected 1 hidden race, got %d", hiddenCount)
	}
}
