package calendar_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Falcon50012/WB-L2/18/internal/calendar"
)

var testDate = time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

func newCal() *calendar.Calendar { return calendar.New() }

func mustCreate(t *testing.T, c *calendar.Calendar, userID int, title string, date time.Time) calendar.Event {
	t.Helper()
	e, err := c.CreateEvent(userID, title, date)
	if err != nil {
		t.Fatalf("CreateEvent: %v", err)
	}
	return e
}

func TestCreateEvent_Success(t *testing.T) {
	c := newCal()
	e, err := c.CreateEvent(1, "Meeting", testDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.ID != 1 {
		t.Errorf("expected ID=1, got %d", e.ID)
	}
	if e.Title != "Meeting" {
		t.Errorf("expected title Meeting, got %s", e.Title)
	}
}

func TestCreateEvent_IDsIncrement(t *testing.T) {
	c := newCal()
	e1 := mustCreate(t, c, 1, "A", testDate)
	e2 := mustCreate(t, c, 1, "B", testDate)
	if e2.ID != e1.ID+1 {
		t.Errorf("expected sequential IDs, got %d and %d", e1.ID, e2.ID)
	}
}

func TestCreateEvent_InvalidUserID(t *testing.T) {
	c := newCal()
	_, err := c.CreateEvent(0, "X", testDate)
	if !errors.Is(err, calendar.ErrInvalidUserID) {
		t.Errorf("expected ErrInvalidUserID, got %v", err)
	}
}

func TestCreateEvent_EmptyTitle(t *testing.T) {
	c := newCal()
	_, err := c.CreateEvent(1, "", testDate)
	if !errors.Is(err, calendar.ErrEmptyTitle) {
		t.Errorf("expected ErrEmptyTitle, got %v", err)
	}
}

func TestUpdateEvent_Success(t *testing.T) {
	c := newCal()
	e := mustCreate(t, c, 1, "Old", testDate)

	newDate := testDate.AddDate(0, 0, 1)
	updated, err := c.UpdateEvent(e.ID, 1, "New", newDate)
	if err != nil {
		t.Fatalf("UpdateEvent: %v", err)
	}
	if updated.Title != "New" {
		t.Errorf("expected New, got %s", updated.Title)
	}
	if !updated.Date.Equal(newDate) {
		t.Errorf("date mismatch: %v vs %v", updated.Date, newDate)
	}
}

func TestUpdateEvent_NotFound(t *testing.T) {
	c := newCal()
	_, err := c.UpdateEvent(999, 1, "X", testDate)
	if !errors.Is(err, calendar.ErrEventNotFound) {
		t.Errorf("expected ErrEventNotFound, got %v", err)
	}
}

func TestUpdateEvent_InvalidUserID(t *testing.T) {
	c := newCal()
	e := mustCreate(t, c, 1, "T", testDate)
	_, err := c.UpdateEvent(e.ID, -1, "T", testDate)
	if !errors.Is(err, calendar.ErrInvalidUserID) {
		t.Errorf("expected ErrInvalidUserID, got %v", err)
	}
}

func TestDeleteEvent_Success(t *testing.T) {
	c := newCal()
	e := mustCreate(t, c, 1, "T", testDate)
	if err := c.DeleteEvent(e.ID); err != nil {
		t.Fatalf("DeleteEvent: %v", err)
	}
}

func TestDeleteEvent_NotFound(t *testing.T) {
	c := newCal()
	err := c.DeleteEvent(42)
	if !errors.Is(err, calendar.ErrEventNotFound) {
		t.Errorf("expected ErrEventNotFound, got %v", err)
	}
}

func TestEventsForDay(t *testing.T) {
	c := newCal()
	mustCreate(t, c, 1, "Today", testDate)
	mustCreate(t, c, 1, "Tomorrow", testDate.AddDate(0, 0, 1))
	mustCreate(t, c, 2, "OtherUser", testDate)

	events, err := c.EventsForDay(1, testDate)
	if err != nil {
		t.Fatalf("EventsForDay: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
	if events[0].Title != "Today" {
		t.Errorf("unexpected title: %s", events[0].Title)
	}
}

func TestEventsForWeek(t *testing.T) {
	c := newCal()

	mustCreate(t, c, 1, "Day0", testDate)
	mustCreate(t, c, 1, "Day6", testDate.AddDate(0, 0, 6))
	mustCreate(t, c, 1, "Day7", testDate.AddDate(0, 0, 7))

	events, err := c.EventsForWeek(1, testDate)
	if err != nil {
		t.Fatalf("EventsForWeek: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

func TestEventsForMonth(t *testing.T) {
	c := newCal()
	june1 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	mustCreate(t, c, 1, "June1", june1)
	mustCreate(t, c, 1, "June30", time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC))
	mustCreate(t, c, 1, "July1", time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC))

	events, err := c.EventsForMonth(1, june1)
	if err != nil {
		t.Fatalf("EventsForMonth: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

func TestConcurrentCreate(t *testing.T) {
	c := newCal()
	done := make(chan struct{})
	const goroutines = 50

	for range goroutines {
		go func() {
			defer func() { done <- struct{}{} }()
			_, _ = c.CreateEvent(1, "concurrent", testDate)
		}()
	}
	for range goroutines {
		<-done
	}
	events, _ := c.EventsForDay(1, testDate)
	if len(events) != goroutines {
		t.Errorf("expected %d events after concurrent creates, got %d", goroutines, len(events))
	}
}
