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

	// Get role from JWT context (set by JWTAuth middleware)
	role, _ := c.Get("role").(string)

	// Fetch user from database
	var user models.Account
	if err := h.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// Return user with role field
	return c.JSON(http.StatusOK, echo.Map{
		"id":       user.ID,
		"username": user.Username,
		"phone":    user.Phone,
		"email":    user.Email,
		"password": user.Password,
		"role":     role,
	})
}

