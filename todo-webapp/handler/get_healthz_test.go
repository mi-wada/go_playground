package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mi-wada/go_playground/todo-webapp/testutils"
)

func TestGetHealthz(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db, err := testutils.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	handler := NewHandler(db)

	if err := handler.GetHealthz(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("invalid status code: %d", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Errorf("invalid response body: %s", rec.Body.String())
	}
}
