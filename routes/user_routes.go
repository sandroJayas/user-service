package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sandroJayas/user-service/controllers"
	"github.com/sandroJayas/user-service/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"net/http"
)

func RegisterUserRoutes(r *gin.Engine, controller *controllers.UserController, db *gorm.DB) {

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/readyz", func(c *gin.Context) {
		db, err := db.DB()
		if err != nil || db.Ping() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db not ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	users := r.Group("/users")
	{
		users.POST("/register", middleware.RateLimitMiddleware(), controller.Register)
		users.POST("/login", middleware.RateLimitMiddleware(), controller.Login)
		users.GET("/me", middleware.AuthMiddleware(), controller.Me)
		users.PUT("/profile", middleware.AuthMiddleware(), controller.UpdateProfile)
		users.DELETE("/delete", middleware.AuthMiddleware(), controller.DeleteUser)

		users.POST("/create-employee", controller.CreateEmployee)
		users.POST("/special", middleware.AuthMiddleware(), middleware.RequireEmployeeRole(), controller.SpecialEmployeeEndpoint)

	}
}
