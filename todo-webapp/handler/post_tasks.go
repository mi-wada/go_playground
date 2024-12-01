package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h handler) PostTasks(c echo.Context) error {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	content := c.FormValue("content")
	status := StatusTodo
	deadlineStr := c.FormValue("deadline")
	var deadline *time.Time
	if deadlineStr != "" {
		parsedDeadline, err := time.Parse(time.RFC3339, deadlineStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, NewError(ErrCodeDeadlineInvalid))
		}
		deadline = &parsedDeadline
	}
	createdAt := time.Now()

	task, err := NewTask(id.String(), content, status, deadline, createdAt)
	switch err {
	case ErrContentRequired:
		return c.JSON(http.StatusBadRequest, NewError(ErrCodeContentRequired))
	case ErrStatusUnknown:
		panic(fmt.Sprintf("status is unknown. err: %v, status: %s", err, status))
	}

	if _, err = h.db.Exec(`
		INSERT INTO
		tasks(id, content, status, deadline, created_at)
		VALUES(?, ?, ?, ?, ?)
	`, task.ID, task.Content, task.Status, task.Deadline, task.CreatedAt); err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError())
	}

	return c.JSON(http.StatusCreated, task)
}

type Task struct {
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	Status    string     `json:"status"`
	Deadline  *time.Time `json:"deadline"`
	CreatedAt time.Time  `json:"created_at"`
}

var (
	ErrContentRequired = errors.New("content is required")
	ErrStatusUnknown   = errors.New("status is unkown")
)

func NewTask(id, content, status string, deadline *time.Time, createdAt time.Time) (*Task, error) {
	if id == "" {
		panic("id is required")
	}
	if content == "" {
		return nil, ErrContentRequired
	}
	if status == "" {
		return nil, ErrStatusUnknown
	}

	return &Task{
		ID:        id,
		Content:   content,
		Status:    status,
		Deadline:  deadline,
		CreatedAt: createdAt,
	}, nil
}

const (
	StatusTodo = "todo"
	StatusDone = "done"
)

type Error struct {
	Code ErrCode `json:"code"`
}

func NewError(code ErrCode) *Error {
	return &Error{Code: code}
}

func InternalServerError() *Error {
	return NewError(ErrCodeInternalServerError)
}

type ErrCode string

const (
	ErrCodeDeadlineInvalid     = ErrCode("deadline_invalid")
	ErrCodeContentRequired     = ErrCode("content_required")
	ErrCodeInternalServerError = ErrCode("internal_server_error")
)
