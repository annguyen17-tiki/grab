package handler

import (
	"net/http"

	"github.com/annguyen17-tiki/grab/internal/model"
	"github.com/annguyen17-tiki/grab/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New() (*gin.Engine, error) {
	if err := loadConfig(); err != nil {
		return nil, err
	}

	if err := initValidator(); err != nil {
		return nil, err
	}

	svc, err := service.New()
	if err != nil {
		return nil, err
	}

	router := gin.Default()
	setCORS(router)

	unauthorized := router.Group("/")
	unauthorized.POST("/register", createAccount(svc))
	unauthorized.POST("/accounts/login", login(svc))

	login := router.Group("/")
	login.Use(authenticate)

	login.GET("/accounts/me", getOwnAccount(svc))
	login.PUT("/accounts/info", updateAccount(svc))

	login.POST("/locations", requireRole(svc, model.RoleDriver), saveLocation(svc))
	login.POST("/bookings", requireRole(svc, model.RoleUser), createBooking(svc))
	login.POST("/bookings/:booking_id/accept", requireRole(svc, model.RoleDriver), acceptBooking(svc))
	login.POST("/bookings/:booking_id/reject", requireRole(svc, model.RoleDriver), rejectBooking(svc))
	login.POST("/bookings/:booking_id/done", requireRole(svc, model.RoleDriver), doneBooking(svc))

	login.GET("/locations/nearest", requireRole(svc, model.RoleUser, model.RoleAdmin), nearestLocations(svc))
	login.GET("/bookings", requireRole(svc, model.RoleUser, model.RoleDriver), searchBookings(svc))
	login.GET("/bookings/:booking_id", getBooking(svc))
	login.GET("/notifications", searchNotifications(svc))

	login.PUT("/notifications/:notification_id/seen", seenNotification(svc))

	admin := router.Group("/admin")
	admin.Use(authenticate, requireRole(svc, model.RoleAdmin))

	admin.GET("/accounts", getAccount(svc))
	admin.GET("/bookings", searchBookingsForAdmin(svc))

	return router, nil
}

func setCORS(engine *gin.Engine) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowMethods(http.MethodOptions)
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	engine.Use(cors.New(corsConfig))
}
