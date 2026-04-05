package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"sports/db"
	"sports/proto/sports"

	_ "github.com/mattn/go-sqlite3"
)

func TestSportsService_ListEvents_VisibleFilter(t *testing.T) {
	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer func() { _ = dbConn.Close() }()

	_, err = dbConn.Exec(`
		CREATE TABLE events (
			id INTEGER PRIMARY KEY,
			name TEXT,
			advertised_start_time DATETIME,
			visible INTEGER
		)
	`)
	if err != nil {
		t.Fatalf("create table error = %v", err)
	}

	start1 := time.Date(2026, 4, 5, 10, 0, 0, 0, time.UTC)
	start2 := time.Date(2026, 4, 5, 11, 0, 0, 0, time.UTC)
	start3 := time.Date(2026, 4, 5, 12, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO events (id, name, advertised_start_time, visible)
		VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?)
	`,
		1, "Visible Event 1", start1, true,
		2, "Hidden Event", start2, false,
		3, "Visible Event 2", start3, true,
	)
	if err != nil {
		t.Fatalf("insert events error = %v", err)
	}

	eventsRepo := db.NewEventsRepo(dbConn)
	service := NewSportsService(eventsRepo)

	// Test filter for visible events only
	visible := true
	request := &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{
			Visible: &visible,
		},
	}

	response, err := service.ListEvents(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}

	if len(response.Events) != 2 {
		t.Fatalf("ListEvents(visible=true) returned %d events, want 2", len(response.Events))
	}

	// Verify both returned events are visible
	for i, event := range response.Events {
		if !event.Visible {
			t.Fatalf("response.Events[%d].Visible = false, want true", i)
		}
	}

	// Verify we got the correct visible events
	expectedIds := []int64{1, 3}
	for i, event := range response.Events {
		if event.Id != expectedIds[i] {
			t.Fatalf("response.Events[%d].Id = %d, want %d", i, event.Id, expectedIds[i])
		}
	}
}

func TestSportsService_ListEvents_HiddenFilter(t *testing.T) {
	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer func() { _ = dbConn.Close() }()

	_, err = dbConn.Exec(`
		CREATE TABLE events (
			id INTEGER PRIMARY KEY,
			name TEXT,
			advertised_start_time DATETIME,
			visible INTEGER
		)
	`)
	if err != nil {
		t.Fatalf("create table error = %v", err)
	}

	start1 := time.Date(2026, 4, 5, 10, 0, 0, 0, time.UTC)
	start2 := time.Date(2026, 4, 5, 11, 0, 0, 0, time.UTC)
	start3 := time.Date(2026, 4, 5, 12, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO events (id, name, advertised_start_time, visible)
		VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?)
	`,
		1, "Visible Event 1", start1, true,
		2, "Hidden Event", start2, false,
		3, "Visible Event 2", start3, true,
	)
	if err != nil {
		t.Fatalf("insert events error = %v", err)
	}

	eventsRepo := db.NewEventsRepo(dbConn)
	service := NewSportsService(eventsRepo)

	// Test filter for hidden events only
	request := &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{
			Visible: new(false),
		},
	}

	response, err := service.ListEvents(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}

	if len(response.Events) != 1 {
		t.Fatalf("ListEvents(visible=false) returned %d events, want 1", len(response.Events))
	}

	// Verify the returned event is hidden
	if response.Events[0].Visible {
		t.Fatalf("response.Events[0].Visible = true, want false")
	}

	// Verify we got the correct hidden event
	if response.Events[0].Id != 2 {
		t.Fatalf("response.Events[0].Id = %d, want 2", response.Events[0].Id)
	}
}

func TestSportsService_ListEvents_NoVisibleFilter(t *testing.T) {
	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer func() { _ = dbConn.Close() }()

	_, err = dbConn.Exec(`
		CREATE TABLE events (
			id INTEGER PRIMARY KEY,
			name TEXT,
			advertised_start_time DATETIME,
			visible INTEGER
		)
	`)
	if err != nil {
		t.Fatalf("create table error = %v", err)
	}

	start1 := time.Date(2026, 4, 5, 10, 0, 0, 0, time.UTC)
	start2 := time.Date(2026, 4, 5, 11, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO events (id, name, advertised_start_time, visible)
		VALUES (?, ?, ?, ?), (?, ?, ?, ?)
	`,
		1, "Visible Event", start1, true,
		2, "Hidden Event", start2, false,
	)
	if err != nil {
		t.Fatalf("insert events error = %v", err)
	}

	eventsRepo := db.NewEventsRepo(dbConn)
	service := NewSportsService(eventsRepo)

	// Test without visibility filter (should return all events)
	request := &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{
			// No Visible field set - should return all events
		},
	}

	response, err := service.ListEvents(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}

	if len(response.Events) != 2 {
		t.Fatalf("ListEvents(no visible filter) returned %d events, want 2", len(response.Events))
	}

	// Verify we got both visible and hidden events
	visibleCount := 0
	hiddenCount := 0
	for _, event := range response.Events {
		if event.Visible {
			visibleCount++
		} else {
			hiddenCount++
		}
	}

	if visibleCount != 1 {
		t.Fatalf("Expected 1 visible event, got %d", visibleCount)
	}

	if hiddenCount != 1 {
		t.Fatalf("Expected 1 hidden event, got %d", hiddenCount)
	}
}

func TestSportsService_ListEvents_OrderByTime(t *testing.T) {
	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer func() { _ = dbConn.Close() }()

	_, err = dbConn.Exec(`
		CREATE TABLE events (
			id INTEGER PRIMARY KEY,
			name TEXT,
			advertised_start_time DATETIME,
			visible INTEGER
		)
	`)
	if err != nil {
		t.Fatalf("create table error = %v", err)
	}

	// Create events with different start times
	start1 := time.Date(2026, 4, 5, 10, 0, 0, 0, time.UTC)
	start2 := time.Date(2026, 4, 5, 12, 0, 0, 0, time.UTC)
	start3 := time.Date(2026, 4, 5, 11, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO events (id, name, advertised_start_time, visible)
		VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?)
	`,
		1, "Early Event", start1, true,
		2, "Late Event", start2, true,
		3, "Mid Event", start3, true,
	)
	if err != nil {
		t.Fatalf("insert events error = %v", err)
	}

	eventsRepo := db.NewEventsRepo(dbConn)
	service := NewSportsService(eventsRepo)

	// Test ascending order (default)
	request := &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{},
	}

	response, err := service.ListEvents(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}

	if len(response.Events) != 3 {
		t.Fatalf("ListEvents() returned %d events, want 3", len(response.Events))
	}

	// Verify ascending order by start time
	expectedIds := []int64{1, 3, 2} // Early, Mid, Late
	for i, event := range response.Events {
		if event.Id != expectedIds[i] {
			t.Fatalf("response.Events[%d].Id = %d, want %d", i, event.Id, expectedIds[i])
		}
	}

	// Test descending order
	request = &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{
			Order: new("time_desc"),
		},
	}

	response, err = service.ListEvents(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}

	// Verify descending order by start time
	expectedIds = []int64{2, 3, 1} // Late, Mid, Early
	for i, event := range response.Events {
		if event.Id != expectedIds[i] {
			t.Fatalf("response.Events[%d].Id = %d, want %d", i, event.Id, expectedIds[i])
		}
	}
}

func TestSportsService_ListEvents_StatusFilter(t *testing.T) {
	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer func() { _ = dbConn.Close() }()

	_, err = dbConn.Exec(`
		CREATE TABLE events (
			id INTEGER PRIMARY KEY,
			name TEXT,
			advertised_start_time DATETIME,
			visible INTEGER
		)
	`)
	if err != nil {
		t.Fatalf("create table error = %v", err)
	}

	// Create events with different start times (some in past for CLOSED, some in future for OPEN)
	pastTime := time.Now().Add(-1 * time.Hour)
	futureTime := time.Now().Add(1 * time.Hour)

	_, err = dbConn.Exec(`
		INSERT INTO events (id, name, advertised_start_time, visible)
		VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?)
	`,
		1, "Closed Event 1", pastTime, true,
		2, "Open Event", futureTime, true,
		3, "Closed Event 2", pastTime.Add(-30*time.Minute), true,
	)
	if err != nil {
		t.Fatalf("insert events error = %v", err)
	}

	eventsRepo := db.NewEventsRepo(dbConn)
	service := NewSportsService(eventsRepo)

	// Test filter for OPEN events only
	status := "OPEN"
	request := &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{
			Status: &status,
		},
	}

	response, err := service.ListEvents(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}

	if len(response.Events) != 1 {
		t.Fatalf("ListEvents(status=OPEN) returned %d events, want 1", len(response.Events))
	}

	// Verify the returned event has OPEN status
	if response.Events[0].Status != "OPEN" {
		t.Fatalf("response.Events[0].Status = %s, want OPEN", response.Events[0].Status)
	}

	// Test filter for CLOSED events only
	status = "CLOSED"
	request = &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{
			Status: &status,
		},
	}

	response, err = service.ListEvents(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}

	if len(response.Events) != 2 {
		t.Fatalf("ListEvents(status=CLOSED) returned %d events, want 2", len(response.Events))
	}

	// Verify all returned events have CLOSED status
	for i, event := range response.Events {
		if event.Status != "CLOSED" {
			t.Fatalf("response.Events[%d].Status = %s, want CLOSED", i, event.Status)
		}
	}
}

func TestSportsService_ListEvents_OrderByName(t *testing.T) {
	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer func() { _ = dbConn.Close() }()

	_, err = dbConn.Exec(`
		CREATE TABLE events (
			id INTEGER PRIMARY KEY,
			name TEXT,
			advertised_start_time DATETIME,
			visible INTEGER
		)
	`)
	if err != nil {
		t.Fatalf("create table error = %v", err)
	}

	start := time.Date(2026, 4, 5, 10, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO events (id, name, advertised_start_time, visible)
		VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?)
	`,
		1, "Zebra Event", start, true,
		2, "Alpha Event", start, true,
		3, "Beta Event", start, true,
	)
	if err != nil {
		t.Fatalf("insert events error = %v", err)
	}

	eventsRepo := db.NewEventsRepo(dbConn)
	service := NewSportsService(eventsRepo)

	// Test ascending order by name
	order := "name_asc"
	request := &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{
			Order: &order,
		},
	}

	response, err := service.ListEvents(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}

	if len(response.Events) != 3 {
		t.Fatalf("ListEvents() returned %d events, want 3", len(response.Events))
	}

	// Verify ascending order by name
	expectedNames := []string{"Alpha Event", "Beta Event", "Zebra Event"}
	for i, event := range response.Events {
		if event.Name != expectedNames[i] {
			t.Fatalf("response.Events[%d].Name = %s, want %s", i, event.Name, expectedNames[i])
		}
	}

	// Test descending order by name
	order = "name_desc"
	request = &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{
			Order: &order,
		},
	}

	response, err = service.ListEvents(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}

	// Verify descending order by name
	expectedNames = []string{"Zebra Event", "Beta Event", "Alpha Event"}
	for i, event := range response.Events {
		if event.Name != expectedNames[i] {
			t.Fatalf("response.Events[%d].Name = %s, want %s", i, event.Name, expectedNames[i])
		}
	}
}

func TestSportsService_ListEvents_InvalidOrder(t *testing.T) {
	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer func() { _ = dbConn.Close() }()

	_, err = dbConn.Exec(`
		CREATE TABLE events (
			id INTEGER PRIMARY KEY,
			name TEXT,
			advertised_start_time DATETIME,
			visible INTEGER
		)
	`)
	if err != nil {
		t.Fatalf("create table error = %v", err)
	}

	start1 := time.Date(2026, 4, 5, 10, 0, 0, 0, time.UTC)
	start2 := time.Date(2026, 4, 5, 12, 0, 0, 0, time.UTC)
	start3 := time.Date(2026, 4, 5, 11, 0, 0, 0, time.UTC)

	_, err = dbConn.Exec(`
		INSERT INTO events (id, name, advertised_start_time, visible)
		VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?)
	`,
		1, "Early Event", start1, true,
		2, "Late Event", start2, true,
		3, "Mid Event", start3, true,
	)
	if err != nil {
		t.Fatalf("insert events error = %v", err)
	}

	eventsRepo := db.NewEventsRepo(dbConn)
	service := NewSportsService(eventsRepo)

	request := &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{
			Order: new("invalid_order"),
		},
	}

	response, err := service.ListEvents(context.TODO(), request)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}

	if len(response.Events) != 3 {
		t.Fatalf("ListEvents() returned %d events, want 3", len(response.Events))
	}

	// Should default to time ascending order
	expectedIds := []int64{1, 3, 2}
	for i, event := range response.Events {
		if event.Id != expectedIds[i] {
			t.Fatalf("response.Events[%d].Id = %d, want %d", i, event.Id, expectedIds[i])
		}
	}
}
