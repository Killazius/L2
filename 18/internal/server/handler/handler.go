package handler

import (
	"18/internal/domain"
	"18/internal/repository"
	"18/internal/service"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	Create(e domain.Event) (domain.Event, error)
	Update(e domain.Event) error
	Delete(userID string, date string, id int64) error
	GetByDate(userID string, date string) ([]domain.Event, error)
	GetByWeek(userID string, date string) ([]domain.Event, error)
	GetByMonth(userID string, date string) ([]domain.Event, error)
}
type Handler struct {
	s Service
}

func New(s Service) *Handler {
	return &Handler{s: s}
}

type CreateEventRequest struct {
	Date   string `json:"date" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
	Text   string `json:"text" binding:"required"`
}

type UpdateEventRequest struct {
	ID     int64  `json:"id" binding:"required"`
	Date   string `json:"date" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
	Text   string `json:"text" binding:"required"`
}

type DeleteEventRequest struct {
	ID     int64  `json:"id" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
	Date   string `json:"date" binding:"required"`
}
type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func (h *Handler) CreateEvent(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			simpleError := simplifyValidationErrors(validationErrors)
			c.JSON(http.StatusBadRequest, Response{Error: simpleError})
			return
		}
		return
	}

	createdEvent, err := h.s.Create(domain.Event{
		Date:   req.Date,
		UserID: req.UserID,
		Text:   req.Text,
	})

	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidDate) {
			status = http.StatusBadRequest
		}
		c.JSON(status, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, Response{Result: createdEvent})
}

func (h *Handler) UpdateEvent(c *gin.Context) {
	var req UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			simpleError := simplifyValidationErrors(validationErrors)
			c.JSON(http.StatusBadRequest, Response{Error: simpleError})
			return
		}
		return
	}

	err := h.s.Update(domain.Event{
		ID:     req.ID,
		Date:   req.Date,
		UserID: req.UserID,
		Text:   req.Text,
	})

	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, repository.ErrNotFound) {
			status = http.StatusServiceUnavailable
		} else if errors.Is(err, service.ErrInvalidDate) {
			status = http.StatusBadRequest
		}
		c.JSON(status, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Result: "event updated successfully"})
}
func (h *Handler) DeleteEvent(c *gin.Context) {
	var req DeleteEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			simpleError := simplifyValidationErrors(validationErrors)
			c.JSON(http.StatusBadRequest, Response{Error: simpleError})
			return
		}
		return
	}

	err := h.s.Delete(req.UserID, req.Date, req.ID)

	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, repository.ErrNotFound) {
			status = http.StatusServiceUnavailable
		} else if errors.Is(err, service.ErrInvalidDate) {
			status = http.StatusBadRequest
		}
		c.JSON(status, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Result: "event deleted"})
}

func (h *Handler) EventsForDay(c *gin.Context) {
	userID, date, ok := getUserIDAndDate(c)
	if !ok {
		return
	}

	events, err := h.s.GetByDate(userID, date)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidDate) {
			status = http.StatusBadRequest
		}
		c.JSON(status, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Result: events})
}

func (h *Handler) EventsForWeek(c *gin.Context) {
	userID, date, ok := getUserIDAndDate(c)
	if !ok {
		return
	}
	events, err := h.s.GetByWeek(userID, date)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidDate) {
			status = http.StatusBadRequest
		}
		c.JSON(status, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Result: events})
}

func (h *Handler) EventsForMonth(c *gin.Context) {
	userID, date, ok := getUserIDAndDate(c)
	if !ok {
		return
	}

	events, err := h.s.GetByMonth(userID, date)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidDate) {
			status = http.StatusBadRequest
		}
		c.JSON(status, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Result: events})
}
func getUserIDAndDate(c *gin.Context) (string, string, bool) {
	userID := c.Query("user_id")
	date := c.Query("date")
	if userID == "" || date == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "user_id and date parameters are required"})
		return "", "", false
	}
	return userID, date, true
}

func simplifyValidationErrors(validationErrors validator.ValidationErrors) string {
	if len(validationErrors) == 0 {
		return "Validation error"
	}

	firstError := validationErrors[0]
	field := firstError.Field()

	switch firstError.Tag() {
	case "required":
		return fmt.Sprintf("Field '%s' is required", field)
	default:
		return fmt.Sprintf("Field '%s' is invalid", field)
	}
}
