package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Falcon50012/WB-L2/18/internal/calendar"
)

const dateLayout = "2006-01-02"

type Service interface {
	CreateEvent(userID int, title string, date time.Time) (calendar.Event, error)
	UpdateEvent(id, userID int, title string, date time.Time) (calendar.Event, error)
	DeleteEvent(id int) error
	EventsForDay(userID int, date time.Time) ([]calendar.Event, error)
	EventsForWeek(userID int, date time.Time) ([]calendar.Event, error)
	EventsForMonth(userID int, date time.Time) ([]calendar.Event, error)
}

type Handler struct {
	svc Service
}

func New(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/create_event", h.CreateEvent)
	mux.HandleFunc("/update_event", h.UpdateEvent)
	mux.HandleFunc("/delete_event", h.DeleteEvent)
	mux.HandleFunc("/events_for_day", h.EventsForDay)
	mux.HandleFunc("/events_for_week", h.EventsForWeek)
	mux.HandleFunc("/events_for_month", h.EventsForMonth)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeResult(w http.ResponseWriter, v any) {
	writeJSON(w, http.StatusOK, map[string]any{"result": v})
}

func writeInputError(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": msg})
}

func writeBusinessError(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": msg})
}

func writeServerError(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": msg})
}

func parseUserID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil || id <= 0 {
		return 0, errors.New("user_id must be a positive integer")
	}
	return id, nil
}

func parseEventID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil || id <= 0 {
		return 0, errors.New("id must be a positive integer")
	}
	return id, nil
}

func parseDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("date is required")
	}
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return time.Time{}, errors.New("date must be in YYYY-MM-DD format")
	}
	return t, nil
}

func parseBody(r *http.Request) (map[string]string, error) {
	ct := r.Header.Get("Content-Type")
	result := make(map[string]string)

	if strings.HasPrefix(ct, "application/json") {
		var raw map[string]any
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			return nil, errors.New("invalid JSON body")
		}
		for k, v := range raw {
			switch val := v.(type) {
			case string:
				result[k] = val
			case float64:
				result[k] = strconv.FormatFloat(val, 'f', -1, 64)
			default:
				return nil, errors.New("unsupported field type")
			}
		}
		return result, nil
	}

	if err := r.ParseForm(); err != nil {
		return nil, errors.New("failed to parse form body")
	}
	for k, vs := range r.Form {
		if len(vs) > 0 {
			result[k] = vs[0]
		}
	}
	return result, nil
}

func isBusinessError(err error) bool {
	return errors.Is(err, calendar.ErrEventNotFound) ||
		errors.Is(err, calendar.ErrEventExists) ||
		errors.Is(err, calendar.ErrInvalidUserID) ||
		errors.Is(err, calendar.ErrEmptyTitle)
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeInputError(w, "method not allowed")
		return
	}

	body, err := parseBody(r)
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	userID, err := parseUserID(body["user_id"])
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	date, err := parseDate(body["date"])
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	title := body["title"]
	if title == "" {
		writeInputError(w, "title is required")
		return
	}

	event, err := h.svc.CreateEvent(userID, title, date)
	if err != nil {
		if isBusinessError(err) {
			writeBusinessError(w, err.Error())
		} else {
			writeServerError(w, err.Error())
		}
		return
	}

	writeResult(w, event)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeInputError(w, "method not allowed")
		return
	}

	body, err := parseBody(r)
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	id, err := parseEventID(body["id"])
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	userID, err := parseUserID(body["user_id"])
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	date, err := parseDate(body["date"])
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	title := body["title"]
	if title == "" {
		writeInputError(w, "title is required")
		return
	}

	event, err := h.svc.UpdateEvent(id, userID, title, date)
	if err != nil {
		if isBusinessError(err) {
			writeBusinessError(w, err.Error())
		} else {
			writeServerError(w, err.Error())
		}
		return
	}

	writeResult(w, event)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeInputError(w, "method not allowed")
		return
	}

	body, err := parseBody(r)
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	id, err := parseEventID(body["id"])
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	if err = h.svc.DeleteEvent(id); err != nil {
		if isBusinessError(err) {
			writeBusinessError(w, err.Error())
		} else {
			writeServerError(w, err.Error())
		}
		return
	}

	writeResult(w, "event deleted")
}

func (h *Handler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeInputError(w, "method not allowed")
		return
	}

	userID, err := parseUserID(r.URL.Query().Get("user_id"))
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	date, err := parseDate(r.URL.Query().Get("date"))
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	events, err := h.svc.EventsForDay(userID, date)
	if err != nil {
		if isBusinessError(err) {
			writeBusinessError(w, err.Error())
		} else {
			writeServerError(w, err.Error())
		}
		return
	}

	writeResult(w, events)
}

func (h *Handler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeInputError(w, "method not allowed")
		return
	}

	userID, err := parseUserID(r.URL.Query().Get("user_id"))
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	date, err := parseDate(r.URL.Query().Get("date"))
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	events, err := h.svc.EventsForWeek(userID, date)
	if err != nil {
		if isBusinessError(err) {
			writeBusinessError(w, err.Error())
		} else {
			writeServerError(w, err.Error())
		}
		return
	}

	writeResult(w, events)
}

func (h *Handler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeInputError(w, "method not allowed")
		return
	}

	userID, err := parseUserID(r.URL.Query().Get("user_id"))
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	date, err := parseDate(r.URL.Query().Get("date"))
	if err != nil {
		writeInputError(w, err.Error())
		return
	}

	events, err := h.svc.EventsForMonth(userID, date)
	if err != nil {
		if isBusinessError(err) {
			writeBusinessError(w, err.Error())
		} else {
			writeServerError(w, err.Error())
		}
		return
	}

	writeResult(w, events)
}
