package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

func (h *SystemHandler) Healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
