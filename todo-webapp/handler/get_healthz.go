package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h handler) GetHealthz(c echo.Context) error {
	err := h.db.Ping()
	if err != nil {
		return c.String(http.StatusInternalServerError, "db error")
	}

	return c.String(http.StatusOK, "ok")
}
