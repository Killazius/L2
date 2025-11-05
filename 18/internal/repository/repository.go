package repository

import (
	"18/internal/domain"
	"errors"
	"sort"
	"sync"
	"time"
)

type Repository struct {
	mu     sync.RWMutex
	store  map[string]map[string][]domain.Event
	nextID int64
}

func New() *Repository {
	return &Repository{
		store:  make(map[string]map[string][]domain.Event),
		nextID: 1,
	}
}

var (
	ErrNotFound = errors.New("event not found")
)

func (r *Repository) Create(e domain.Event) (domain.Event, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	e.ID = r.nextID
	r.nextID++

	if r.store[e.UserID] == nil {
		r.store[e.UserID] = make(map[string][]domain.Event)
	}

	r.store[e.UserID][e.Date] = append(r.store[e.UserID][e.Date], e)

	return e, nil
}

func (r *Repository) Update(e domain.Event) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	userEvents, userExists := r.store[e.UserID]
	if !userExists {
		return ErrNotFound
	}

	events, dateExists := userEvents[e.Date]
	if !dateExists {
		return ErrNotFound
	}

	found := false
	for i, event := range events {
		if event.ID == e.ID {
			events[i] = e
			found = true
			break
		}
	}

	if !found {
		return ErrNotFound
	}

	r.store[e.UserID][e.Date] = events
	return nil
}

func (r *Repository) Delete(userID, date string, id int64) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	userEvents, userExists := r.store[userID]
	if !userExists {
		return ErrNotFound
	}

	events, dateExists := userEvents[date]
	if !dateExists {
		return ErrNotFound
	}

	found := false
	newEvents := make([]domain.Event, 0, len(events)-1)
	for _, event := range events {
		if event.ID != id {
			newEvents = append(newEvents, event)
		} else {
			found = true
		}
	}

	if !found {
		return ErrNotFound
	}

	if len(newEvents) == 0 {
		delete(userEvents, date)
		if len(userEvents) == 0 {
			delete(r.store, userID)
		}
	} else {
		r.store[userID][date] = newEvents
	}

	return nil
}

func (r *Repository) GetByDate(userID, date string) ([]domain.Event, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	userEvents, userExists := r.store[userID]
	if !userExists {
		return []domain.Event{}, nil
	}

	events, dateExists := userEvents[date]
	if !dateExists {
		return []domain.Event{}, nil
	}

	out := make([]domain.Event, len(events))
	copy(out, events)
	return out, nil
}

func (r *Repository) GetByWeek(userID, date string) ([]domain.Event, error) {

	startDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	userEvents, userExists := r.store[userID]
	if !userExists {
		return []domain.Event{}, nil
	}

	var events []domain.Event
	for i := 0; i < 7; i++ {
		currentDate := startDate.AddDate(0, 0, i).Format("2006-01-02")
		if dayEvents, exists := userEvents[currentDate]; exists {
			events = append(events, dayEvents...)
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].ID < events[j].ID
	})

	return events, nil
}

func (r *Repository) GetByMonth(userID, date string) ([]domain.Event, error) {

	startDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	userEvents, userExists := r.store[userID]
	if !userExists {
		return []domain.Event{}, nil
	}

	var events []domain.Event
	year, month, _ := startDate.Date()
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, startDate.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	for d := firstOfMonth; !d.After(lastOfMonth); d = d.AddDate(0, 0, 1) {
		currentDate := d.Format("2006-01-02")
		if dayEvents, exists := userEvents[currentDate]; exists {
			events = append(events, dayEvents...)
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].ID < events[j].ID
	})

	return events, nil
}
