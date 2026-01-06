package handlers

import (
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/utils"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB     *gorm.DB
	Secret string
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginReq

	// Read + parse to JSON body
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	// Read account from DB by username
	var acc models.Account
	if err := h.DB.Where("username = ?", req.Username).First(&acc).Error; err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "wrong username or password")
	}

	// Check if password hash match
	if !utils.VerifyPassword(acc.Password, req.Password) {
		return echo.NewHTTPError(http.StatusUnauthorized, "wrong username or password")
	}

	role := "member"
	if acc.Username == "admin" {
		role = "admin"
	}

	// flow: check admin permission -> frontend route theo role
	// JWT Claim
	claims := middleware.JwtClaims{
		UserID: uint(acc.ID),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.Secret))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "sign token failed")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": tokenStr,
		"role":  role,
	})
}
