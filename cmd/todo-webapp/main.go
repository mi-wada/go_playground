package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mi-wada/go_playground/todo-webapp/handler"
)

func main() {
	e := echo.New()

	e.GET("/healthz", handler.GetHealthz)

	e.Start(":8080")
}
