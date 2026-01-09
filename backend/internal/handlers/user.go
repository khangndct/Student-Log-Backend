package handlers

import (
	"backend/internal/models"
	"backend/internal/utils"
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

type UpdateAccountRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Phone    *int64  `json:"phone"`
	Password *string `json:"password"`
}

func (h *UserHandler) Update(c echo.Context) error {
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

	// Bind update request
	var req UpdateAccountRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	// Update only provided fields
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Password != nil {
		hash, err := utils.HashPassword(*req.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "hash failed")
		}
		user.Password = hash
	}

	// Save updated user
	if err := h.DB.Save(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "update failed")
	}

	// Return updated user with role field (same structure as GET /api/user)
	return c.JSON(http.StatusOK, echo.Map{
		"id":       user.ID,
		"username": user.Username,
		"phone":    user.Phone,
		"email":    user.Email,
		"password": user.Password,
		"role":     role,
	})
}

