package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mi-wada/go_playground/todo-webapp/testutils"
)

func TestOkPostTasks(t *testing.T) {
	type testCase struct {
		name     string
		content  string
		deadline string
	}

	for _, tc := range []testCase{
		{
			name:     "with deadline",
			content:  "content",
			deadline: "2021-01-01T00:00:00Z",
		},
		{
			name:     "no deadline",
			content:  "content",
			deadline: "",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				strings.NewReader("content="+tc.content+"&deadline="+tc.deadline),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			db, err := testutils.InitDB()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			handler := NewHandler(db)

			if err := handler.PostTasks(c); err != nil {
				t.Fatal(err)
			}
			if rec.Code != http.StatusCreated {
				t.Errorf("invalid status code: %d", rec.Code)
			}
			var task Task
			err = json.Unmarshal(rec.Body.Bytes(), &task)
			if err != nil {
				t.Fatal(err)
			}
			if !(len(task.ID) > 0 &&
				task.Content == tc.content &&
				task.Status == StatusTodo &&
				(tc.deadline == "" || task.Deadline.String() == "2021-01-01 00:00:00 +0000 UTC")) {
				t.Errorf("invalid response body: %v", task)
			}
		})
	}
}

func TestErrPostTasks(t *testing.T) {
	type testCase struct {
		name     string
		content  string
		deadline string
		errCode  string
	}

	for tc := range slices.Values(
		[]*testCase{
			{
				name:     "content is required",
				content:  "",
				deadline: "2021-01-01T00:00:00Z",
				errCode:  "content_required",
			},
			{
				name:     "deadline is invalid",
				content:  "content",
				deadline: "invalid",
				errCode:  "deadline_invalid",
			},
		},
	) {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				strings.NewReader("content="+tc.content+"&deadline="+tc.deadline),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			db, err := testutils.InitDB()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			handler := NewHandler(db)

			if err := handler.PostTasks(c); err != nil {
				t.Fatal(err)
			}
			if rec.Code != http.StatusBadRequest {
				t.Errorf("invalid status code: %d", rec.Code)
			}
			if !strings.Contains(rec.Body.String(), tc.errCode) {
				t.Errorf("invalid response body: %s", rec.Body.String())
			}
		})
	}
}
