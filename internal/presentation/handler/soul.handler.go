package handler

import (
	"itemsim-server/internal/application"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SoulHandler struct {
	soulService application.SoulService
}

func NewSoulHandler(soulService application.SoulService) *SoulHandler {
	return &SoulHandler{
		soulService: soulService,
	}
}

func (h *SoulHandler) GetAllDataAsJson(c echo.Context) error {
	result, err := h.soulService.GetAllDataAsJson()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}
