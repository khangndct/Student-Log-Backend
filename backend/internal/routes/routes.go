package routes

import (
	"backend/internal/handlers"
	appmw "backend/internal/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Register(e *echo.Echo, db *gorm.DB, jwtSecret string) {
	// Auth
	auth := &handlers.AuthHandler{DB: db, Secret: jwtSecret}
	e.POST("/api/auth/login", auth.Login)

	// Protected
	api := e.Group("/api", appmw.JWTAuth(jwtSecret))

	// Admin
	adminAccounts := &handlers.AdminAccountsHandler{DB: db}
	adminLogs := &handlers.AdminLogsHandler{DB: db}

	admin := api.Group("/admin", appmw.RequireRole("admin"))
	admin.GET("/accounts", adminAccounts.List)
	admin.POST("/accounts", adminAccounts.Create)
	admin.DELETE("/accounts/:id", adminAccounts.Delete)

	admin.GET("/log-heads", adminLogs.ListLogHeads)
	admin.POST("/log-heads", adminLogs.CreateLogHead)
	admin.PUT("/log-heads/:id", adminLogs.UpdateLogHead)
	admin.PATCH("/log-heads/:id", adminLogs.UpdateLogHead)
	admin.DELETE("/log-heads/:id", adminLogs.DeleteLogHead)

	memberLogs := &handlers.MemberLogsHandler{DB: db}
	api.GET("/log-heads", memberLogs.ListLogHeads)
	api.GET("/log-heads/writable", memberLogs.ListWritableLogHeads)
	api.POST("/log-contents", memberLogs.CreateLogContent)
	api.PUT("/log-contents/:id", memberLogs.UpdateLogContent)
	api.PATCH("/log-contents/:id", memberLogs.UpdateLogContent)
	api.DELETE("/log-contents/:id", memberLogs.DeleteLogContent)

	// Members search (admin only)
	membersSearch := api.Group("/members", appmw.RequireRole("admin"))
	membersSearch.GET("/search", adminAccounts.SearchMembers)

	// User
	user := &handlers.UserHandler{DB: db}
	api.GET("/user", user.Get)
}
