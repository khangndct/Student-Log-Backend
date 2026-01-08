package handlers

import (
	"backend/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func (h *UserHandler) Get(c echo.Context) error {
	// Get user_id from JWT context (set by JWTAuth middleware)
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	// Fetch user from database
	var user models.Account
	if err := h.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, user)
}

