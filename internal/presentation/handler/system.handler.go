package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

func (h *SystemHandler) Healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
