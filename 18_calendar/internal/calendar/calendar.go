package calendar

import (
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrEventExists   = errors.New("event with this id already exists")
	ErrInvalidUserID = errors.New("user_id must be greater than zero")
	ErrEmptyTitle    = errors.New("event title must not be empty")
)

type Event struct {
	ID     int       `json:"id"`
	UserID int       `json:"user_id"`
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
}

type Calendar struct {
	mu     sync.RWMutex
	events map[int]Event
	nextID int
}

func New() *Calendar {
	return &Calendar{
		events: make(map[int]Event),
		nextID: 1,
	}
}

func (c *Calendar) CreateEvent(userID int, title string, date time.Time) (Event, error) {
	if err := validateUserID(userID); err != nil {
		return Event{}, err
	}
	if err := validateTitle(title); err != nil {
		return Event{}, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	e := Event{
		ID:     c.nextID,
		UserID: userID,
		Title:  title,
		Date:   date.Truncate(24 * time.Hour),
	}
	c.events[c.nextID] = e
	c.nextID++
	return e, nil
}

func (c *Calendar) UpdateEvent(id, userID int, title string, date time.Time) (Event, error) {
	if err := validateUserID(userID); err != nil {
		return Event{}, err
	}
	if err := validateTitle(title); err != nil {
		return Event{}, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.events[id]; !ok {
		return Event{}, ErrEventNotFound
	}

	e := Event{
		ID:     id,
		UserID: userID,
		Title:  title,
		Date:   date.Truncate(24 * time.Hour),
	}
	c.events[id] = e
	return e, nil
}

func (c *Calendar) DeleteEvent(id int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.events[id]; !ok {
		return ErrEventNotFound
	}
	delete(c.events, id)
	return nil
}

func (c *Calendar) EventsForDay(userID int, date time.Time) ([]Event, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}
	day := date.Truncate(24 * time.Hour)

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.filter(func(e Event) bool {
		return e.UserID == userID && e.Date.Equal(day)
	}), nil
}

func (c *Calendar) EventsForWeek(userID int, date time.Time) ([]Event, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}
	start := date.Truncate(24 * time.Hour)
	end := start.AddDate(0, 0, 7)

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.filter(func(e Event) bool {
		return e.UserID == userID &&
			(e.Date.Equal(start) || e.Date.After(start)) &&
			e.Date.Before(end)
	}), nil
}

func (c *Calendar) EventsForMonth(userID int, date time.Time) ([]Event, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}
	y, m, _ := date.Date()
	start := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.filter(func(e Event) bool {
		return e.UserID == userID &&
			(e.Date.Equal(start) || e.Date.After(start)) &&
			e.Date.Before(end)
	}), nil
}

func (c *Calendar) filter(predicate func(Event) bool) []Event {
	result := make([]Event, 0)
	for _, e := range c.events {
		if predicate(e) {
			result = append(result, e)
		}
	}
	return result
}

func validateUserID(id int) error {
	if id <= 0 {
		return ErrInvalidUserID
	}
	return nil
}

func validateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrEmptyTitle
	}
	return nil
}
