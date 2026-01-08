package handlers

import (
	"backend/internal/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AdminLogsHandler struct {
	DB *gorm.DB
}

func (h *AdminLogsHandler) ListLogHeads(c echo.Context) error {
	var heads []models.LogHead
	if err := h.DB.Preload("LogContents").Find(&heads).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "db error")
	}
	return c.JSON(http.StatusOK, heads)
}

type CreateLogHeadReq struct {
	Subject      string        `json:"subject"`
	StartDate    time.Time     `json:"start_date"`
	EndDate      time.Time     `json:"end_date"`
	WriterIDList pq.Int64Array `json:"writer_id_list"`
	OwnerID      uint          `json:"owner_id"`
}

func (h *AdminLogsHandler) CreateLogHead(c echo.Context) error {
	var req CreateLogHeadReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	head := models.LogHead{
		Subject:      req.Subject,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		WriterIDList: req.WriterIDList,
		OwnerID:      req.OwnerID,
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

type UpdateLogHeadReq struct {
	Subject      *string        `json:"subject"`
	StartDate    *time.Time     `json:"start_date"`
	EndDate      *time.Time     `json:"end_date"`
	WriterIDList *pq.Int64Array `json:"writer_id_list"`
	OwnerID      *uint          `json:"owner_id"`
}

func (h *AdminLogsHandler) UpdateLogHead(c echo.Context) error {
	id := c.Param("id")
	
	var head models.LogHead
	if err := h.DB.Preload("LogContents").First(&head, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "log head not found")
	}

	var req UpdateLogHeadReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	// Update only provided fields
	if req.Subject != nil {
		head.Subject = *req.Subject
	}
	if req.StartDate != nil {
		head.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		head.EndDate = *req.EndDate
	}
	if req.WriterIDList != nil {
		head.WriterIDList = *req.WriterIDList
	}
	if req.OwnerID != nil {
		head.OwnerID = *req.OwnerID
	}

	if err := h.DB.Save(&head).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "update failed")
	}

	return c.JSON(http.StatusOK, head)
}