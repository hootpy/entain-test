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

func TestRacingService_ListRaces_OrderByTime(t *testing.T) {
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

	// Create races with different start times
	start1 := time.Date(2026, 3, 31, 10, 0, 0, 0, time.UTC)
	start2 := time.Date(2026, 3, 31, 12, 0, 0, 0, time.UTC)
	start3 := time.Date(2026, 3, 31, 11, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO races (id, meeting_id, name, number, visible, advertised_start_time)
		VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)
	`,
		1, 100, "Early Race", 1, true, start1,
		2, 200, "Late Race", 2, true, start2,
		3, 300, "Mid Race", 3, true, start3,
	)
	if err != nil {
		t.Fatalf("insert races error = %v", err)
	}

	racesRepo := db.NewRacesRepo(dbConn)
	service := NewRacingService(racesRepo)

	// Test ascending order (default)
	request := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{},
	}

	response, err := service.ListRaces(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListRaces() error = %v", err)
	}

	if len(response.Races) != 3 {
		t.Fatalf("ListRaces() returned %d races, want 3", len(response.Races))
	}

	// Verify ascending order by start time
	expectedIds := []int64{1, 3, 2} // Early, Mid, Late
	for i, race := range response.Races {
		if race.Id != expectedIds[i] {
			t.Fatalf("response.Races[%d].Id = %d, want %d", i, race.Id, expectedIds[i])
		}
	}

	// Test descending order
	order := "time_desc"
	request = &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{
			Order: &order,
		},
	}

	response, err = service.ListRaces(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListRaces() error = %v", err)
	}

	// Verify descending order by start time
	expectedIds = []int64{2, 3, 1} // Late, Mid, Early
	for i, race := range response.Races {
		if race.Id != expectedIds[i] {
			t.Fatalf("response.Races[%d].Id = %d, want %d", i, race.Id, expectedIds[i])
		}
	}
}

func TestRacingService_ListRaces_OrderByName(t *testing.T) {
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

	start := time.Date(2026, 3, 31, 10, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO races (id, meeting_id, name, number, visible, advertised_start_time)
		VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)
	`,
		1, 100, "Zebra Race", 1, true, start,
		2, 200, "Alpha Race", 2, true, start,
		3, 300, "Beta Race", 3, true, start,
	)
	if err != nil {
		t.Fatalf("insert races error = %v", err)
	}

	racesRepo := db.NewRacesRepo(dbConn)
	service := NewRacingService(racesRepo)

	// Test ascending order by name
	order := "name_asc"
	request := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{
			Order: &order,
		},
	}

	response, err := service.ListRaces(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListRaces() error = %v", err)
	}

	if len(response.Races) != 3 {
		t.Fatalf("ListRaces() returned %d races, want 3", len(response.Races))
	}

	// Verify ascending order by name
	expectedNames := []string{"Alpha Race", "Beta Race", "Zebra Race"}
	for i, race := range response.Races {
		if race.Name != expectedNames[i] {
			t.Fatalf("response.Races[%d].Name = %s, want %s", i, race.Name, expectedNames[i])
		}
	}

	// Test descending order by name
	order = "name_desc"
	request = &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{
			Order: &order,
		},
	}

	response, err = service.ListRaces(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListRaces() error = %v", err)
	}

	// Verify descending order by name
	expectedNames = []string{"Zebra Race", "Beta Race", "Alpha Race"}
	for i, race := range response.Races {
		if race.Name != expectedNames[i] {
			t.Fatalf("response.Races[%d].Name = %s, want %s", i, race.Name, expectedNames[i])
		}
	}
}

func TestRacingService_ListRaces_OrderByMeetingId(t *testing.T) {
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

	start := time.Date(2026, 3, 31, 10, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO races (id, meeting_id, name, number, visible, advertised_start_time)
		VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)
	`,
		1, 300, "Race 1", 1, true, start,
		2, 100, "Race 2", 2, true, start,
		3, 200, "Race 3", 3, true, start,
	)
	if err != nil {
		t.Fatalf("insert races error = %v", err)
	}

	racesRepo := db.NewRacesRepo(dbConn)
	service := NewRacingService(racesRepo)

	// Test ascending order by meeting ID
	order := "meeting_asc"
	request := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{
			Order: &order,
		},
	}

	response, err := service.ListRaces(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListRaces() error = %v", err)
	}

	if len(response.Races) != 3 {
		t.Fatalf("ListRaces() returned %d races, want 3", len(response.Races))
	}

	// Verify ascending order by meeting ID
	expectedIds := []int64{2, 3, 1} // Meeting 100, 200, 300
	for i, race := range response.Races {
		if race.Id != expectedIds[i] {
			t.Fatalf("response.Races[%d].Id = %d, want %d", i, race.Id, expectedIds[i])
		}
	}

	// Test descending order by meeting ID
	order = "meeting_desc"
	request = &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{
			Order: &order,
		},
	}

	response, err = service.ListRaces(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListRaces() error = %v", err)
	}

	// Verify descending order by meeting ID
	expectedIds = []int64{1, 3, 2} // Meeting 300, 200, 100
	for i, race := range response.Races {
		if race.Id != expectedIds[i] {
			t.Fatalf("response.Races[%d].Id = %d, want %d", i, race.Id, expectedIds[i])
		}
	}
}

func TestRacingService_ListRaces_InvalidOrder(t *testing.T) {
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
	start2 := time.Date(2026, 3, 31, 12, 0, 0, 0, time.UTC)
	start3 := time.Date(2026, 3, 31, 11, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO races (id, meeting_id, name, number, visible, advertised_start_time)
		VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)
	`,
		1, 100, "Early Race", 1, true, start1,
		2, 200, "Late Race", 2, true, start2,
		3, 300, "Mid Race", 3, true, start3,
	)
	if err != nil {
		t.Fatalf("insert races error = %v", err)
	}

	racesRepo := db.NewRacesRepo(dbConn)
	service := NewRacingService(racesRepo)

	order := "invalid_order"
	request := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{
			Order: &order,
		},
	}

	response, err := service.ListRaces(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListRaces() error = %v", err)
	}

	if len(response.Races) != 3 {
		t.Fatalf("ListRaces() returned %d races, want 3", len(response.Races))
	}

	expectedIds := []int64{1, 3, 2}
	for i, race := range response.Races {
		if race.Id != expectedIds[i] {
			t.Fatalf("response.Races[%d].Id = %d, want %d", i, race.Id, expectedIds[i])
		}
	}
}

func TestRacingService_GetRace_Found(t *testing.T) {
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

	start := time.Date(2026, 3, 31, 10, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO races (id, meeting_id, name, number, visible, advertised_start_time)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		1, 100, "Race One", 1, true, start,
	)
	if err != nil {
		t.Fatalf("insert races error = %v", err)
	}

	racesRepo := db.NewRacesRepo(dbConn)
	service := NewRacingService(racesRepo)

	response, err := service.GetRace(context.TODO(), &racing.GetRaceRequest{Id: 1})
	if err != nil {
		t.Fatalf("GetRace() error = %v", err)
	}

	if response == nil {
		t.Fatal("GetRace() response is nil")
	}

	if response.Race == nil {
		t.Fatal("GetRace().Race is nil")
	}

	if response.Race.Id != 1 {
		t.Fatalf("GetRace().Race.Id = %d, want 1", response.Race.Id)
	}

	if response.Race.MeetingId != 100 {
		t.Fatalf("GetRace().Race.MeetingId = %d, want 100", response.Race.MeetingId)
	}

	if response.Race.Name != "Race One" {
		t.Fatalf("GetRace().Race.Name = %q, want %q", response.Race.Name, "Race One")
	}

	if response.Race.Number != 1 {
		t.Fatalf("GetRace().Race.Number = %d, want 1", response.Race.Number)
	}

	if !response.Race.Visible {
		t.Fatal("GetRace().Race.Visible = false, want true")
	}

	if response.Race.AdvertisedStartTime == nil {
		t.Fatal("GetRace().Race.AdvertisedStartTime is nil")
	}
}

func TestRacingService_GetRace_NotFound(t *testing.T) {
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

	start := time.Date(2026, 3, 31, 10, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO races (id, meeting_id, name, number, visible, advertised_start_time)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		1, 100, "Race One", 1, true, start,
	)
	if err != nil {
		t.Fatalf("insert races error = %v", err)
	}

	racesRepo := db.NewRacesRepo(dbConn)
	service := NewRacingService(racesRepo)

	response, err := service.GetRace(context.TODO(), &racing.GetRaceRequest{Id: 999})
	if err != nil {
		t.Fatalf("GetRace() error = %v", err)
	}

	if response == nil {
		t.Fatal("GetRace() response is nil")
	}

	if response.Race != nil {
		t.Fatalf("GetRace().Race = %#v, want nil", response.Race)
	}
}
