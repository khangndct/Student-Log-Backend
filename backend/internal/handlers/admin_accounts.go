package handlers

import (
	"backend/internal/models"
	"backend/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AdminAccountsHandler struct {
	DB *gorm.DB
}

func (h *AdminAccountsHandler) List(c echo.Context) error {
	var accounts []models.Account
	if err := h.DB.Find(&accounts).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "db error")
	}
	return c.JSON(http.StatusOK, accounts)
}

type CreateAccountReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    int64  `json:"phone"`
	Password string `json:"password"`
}

func (h *AdminAccountsHandler) Create(c echo.Context) error {
	var req CreateAccountReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	// flow: hashing password
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "hash failed")
	}

	acc := models.Account{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hash,
	}
	if err := h.DB.Create(&acc).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "create failed (duplicate username?)")
	}

	return c.JSON(http.StatusCreated, acc)
}

func (h *AdminAccountsHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.DB.Delete(&models.Account{}, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "delete failed")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *AdminAccountsHandler) SearchMembers(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "query parameter 'q' is required")
	}

	var accounts []models.Account
	searchPattern := "%" + query + "%"
	
	// Search in username, email, or phone
	if err := h.DB.Where(
		"username ILIKE ? OR email ILIKE ? OR CAST(phone AS TEXT) LIKE ?",
		searchPattern,
		searchPattern,
		searchPattern,
	).Find(&accounts).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "db error")
	}

	return c.JSON(http.StatusOK, accounts)
}