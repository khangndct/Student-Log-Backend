package handlers

import (
	"backend/internal/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type MemberLogsHandler struct {
	DB *gorm.DB
}

func (h *MemberLogsHandler) ListLogHeads(c echo.Context) error {
	var heads []models.LogHead
	if err := h.DB.Preload("LogContents").Find(&heads).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "db error")
	}
	return c.JSON(http.StatusOK, heads)
}

// flow: filter writable log heads by WriterIDList
func (h *MemberLogsHandler) ListWritableLogHeads(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	role := c.Get("role").(string)

	var heads []models.LogHead
	q := h.DB.Model(&models.LogHead{})

	if role == "admin" {
		if err := q.Preload("LogContents").Find(&heads).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "db error")
		}
		return c.JSON(http.StatusOK, heads)
	}

	if err := q.Preload("LogContents").Where("? = ANY(writer_id_list)", int64(userID)).Find(&heads).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "db error")
	}

	return c.JSON(http.StatusOK, heads)
}

type CreateLogContentReq struct {
	LogHeadID uint      `json:"log_head_id"`
	Content   string    `json:"content"`
	Date      time.Time `json:"date"` // frontend gửi ISO time
}

// flow: chọn log -> lấy log_head_id -> nhập content + timestamp -> ghi db (kèm user_id)
func (h *MemberLogsHandler) CreateLogContent(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	role := c.Get("role").(string)

	var req CreateLogContentReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	var head models.LogHead
	if err := h.DB.First(&head, req.LogHeadID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "log head not found")
	}

	// kiểm tra permission
	canWrite := role == "admin"
	if !canWrite {
		for _, writerID := range head.WriterIDList {
			if writerID == int64(userID) {
				canWrite = true
				break
			}
		}
	}
	if !canWrite {
		return echo.NewHTTPError(http.StatusForbidden, "no write permission for this log")
	}

	lc := models.LogContent{
		LogHeadID: head.ID,
		WriterID:  userID,
		Content:   req.Content,
		Date:      req.Date,
	}
	if err := h.DB.Create(&lc).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "create failed")
	}

	return c.JSON(http.StatusCreated, lc)
}

type UpdateLogContentReq struct {
	Content *string    `json:"content"`
	Date    *time.Time `json:"date"`
}

func (h *MemberLogsHandler) UpdateLogContent(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("user_id").(uint)
	role := c.Get("role").(string)

	var content models.LogContent
	if err := h.DB.First(&content, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "log content not found")
	}

	// Check permission: admin or the original writer
	if role != "admin" && content.WriterID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "no permission to update this log content")
	}

	var req UpdateLogContentReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	// Update only provided fields
	if req.Content != nil {
		content.Content = *req.Content
	}
	if req.Date != nil {
		content.Date = *req.Date
	}

	if err := h.DB.Save(&content).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "update failed")
	}

	return c.JSON(http.StatusOK, content)
}

func (h *MemberLogsHandler) DeleteLogContent(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("user_id").(uint)
	role := c.Get("role").(string)

	var content models.LogContent
	if err := h.DB.First(&content, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "log content not found")
	}

	// Check permission: admin or the original writer
	if role != "admin" && content.WriterID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "no permission to delete this log content")
	}

	if err := h.DB.Delete(&content).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "delete failed")
	}

	return c.NoContent(http.StatusNoContent)
}