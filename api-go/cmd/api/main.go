package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"os"
	"github.com/gin-gonic/gin"

	"github.com/Issarat/food-delivery-swagger-api-go/api-go/docs"
	"github.com/Issarat/food-delivery-swagger-api-go/api-go/internal/config"
	"github.com/Issarat/food-delivery-swagger-api-go/api-go/internal/postgres"
	"github.com/Issarat/food-delivery-swagger-api-go/api-go/internal/middleware"
	"github.com/Issarat/food-delivery-swagger-api-go/api-go/internal/auth"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	db, err := postgres.Open(ctx, cfg.DatabaseURL)
	cancel()
	if err != nil {
		return err
	}
	defer db.Close()
	tokens, err := auth.NewTokenManager(os.Getenv("JWT_SECRET"))
	if err != nil {
		return err
	}
	router := gin.Default()
	authRepo := auth.NewPostgresAuthRepo(db)
	authService := auth.NewService(authRepo, tokens)
	
	requireJWT := middleware.RequireJWT(tokens, authService)

	docs.Register(router)
	router.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "database unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	
	authHandler := auth.NewHandler(authService)

	authGroup := router.Group("/auth")
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/register", authHandler.Register)

	authGroup.POST("/logout", requireJWT, authHandler.Logout)
	authGroup.POST("/change-password", requireJWT, authHandler.ChangePassword)

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return server.ListenAndServe()
}
