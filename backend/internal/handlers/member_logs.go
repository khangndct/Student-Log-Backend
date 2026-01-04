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
	if err := h.DB.Find(&heads).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "db error")
	}
	return c.JSON(http.StatusOK, heads)
}

// flow: check write permission từng log -> trả danh sách writable
func (h *MemberLogsHandler) ListWritableLogHeads(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	role := c.Get("role").(string)

	var heads []models.LogHead
	q := h.DB.Model(&models.LogHead{})

	// member chỉ được ghi:
	// - write_scope = all
	// - write_scope = owner và owner_id = mình
	// admin thì (tuỳ bạn) có thể ghi hết; nhưng flow member nên cứ giữ logic rõ ràng
	if role == "admin" {
		if err := q.Find(&heads).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "db error")
		}
		return c.JSON(http.StatusOK, heads)
	}

	if err := q.Where("write_scope = ? OR (write_scope = ? AND owner_id = ?)", "all", "owner", userID).
		Find(&heads).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "db error")
	}
	return c.JSON(http.StatusOK, heads)
}

type CreateLogContentReq struct {
	LogHeadID uint      `json:"log_head_id"`
	Content   string    `json:"content"`
	LogTime   time.Time `json:"log_time"` // frontend gửi ISO time
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
	canWrite := false
	if role == "admin" {
		canWrite = true
	} else {
		switch head.WriteScope {
		case "all":
			canWrite = true
		case "owner":
			canWrite = head.OwnerID == userID
		case "admin":
			canWrite = false
		}
	}
	if !canWrite {
		return echo.NewHTTPError(http.StatusForbidden, "no write permission for this log")
	}

	lc := models.LogContent{
		LogHeadID: head.ID,
		UserID:    userID,
		Content:   req.Content,
		LogTime:   req.LogTime,
	}
	if err := h.DB.Create(&lc).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "create failed")
	}

	return c.JSON(http.StatusCreated, lc)
}
