package handlers

import (
	"backend/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AdminLogsHandler struct {
	DB *gorm.DB
}

func (h *AdminLogsHandler) ListLogHeads(c echo.Context) error {
	var heads []models.LogHead
	if err := h.DB.Find(&heads).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "db error")
	}
	return c.JSON(http.StatusOK, heads)
}

type CreateLogHeadReq struct {
	Title      string `json:"title"`
	WriteScope string `json:"write_scope"` // "all" | "owner" | "admin"
	OwnerID    uint   `json:"owner_id"`
}

func (h *AdminLogsHandler) CreateLogHead(c echo.Context) error {
	var req CreateLogHeadReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	if req.WriteScope != "all" && req.WriteScope != "owner" && req.WriteScope != "admin" {
		return echo.NewHTTPError(http.StatusBadRequest, "write_scope must be all|owner|admin")
	}

	head := models.LogHead{
		Title:      req.Title,
		WriteScope: req.WriteScope,
		OwnerID:    req.OwnerID,
	}
	if err := h.DB.Create(&head).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "create failed")
	}
	return c.JSON(http.StatusCreated, head)
}

func (h *AdminLogsHandler) DeleteLogHead(c echo.Context) error {
	id := c.Param("id")
	// Cascade auto delete LogContent
	if err := h.DB.Delete(&models.LogHead{}, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "delete failed")
	}
	return c.NoContent(http.StatusNoContent)
}
