package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sandroJayas/user-service/config"
	"github.com/sandroJayas/user-service/controllers"
	"github.com/sandroJayas/user-service/infrastructure/repository"
	"github.com/sandroJayas/user-service/routes"
	"github.com/sandroJayas/user-service/usecase"
	"github.com/sandroJayas/user-service/utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/sandroJayas/user-service/docs"
)

func main() {
	utils.InitLogger()
	defer utils.Logger.Sync()

	config.LoadEnv()
	db := config.ConnectDB()

	userRepo := repository.NewGormUserRepository(db)
	userService := usecase.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	shutdown := utils.InitTracer()
	defer shutdown(context.Background())

	r := gin.Default()
	routes.RegisterUserRoutes(r, userController, db)
	r.Use(otelgin.Middleware("user-service"))

	//graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		utils.Logger.Info("🚀 Server is running on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.Logger.Info("🛑 Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.Logger.Error("failed to gracefully shutdown", zap.Error(err))
	} else {
		utils.Logger.Info("✅ Server shutdown completed")
	}
}
