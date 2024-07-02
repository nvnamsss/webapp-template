package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"webapp-template/src/configs"
	"webapp-template/src/logger"
	"webapp-template/src/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const AppExitCode = 99

func main() {

	corsConfig := cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})

	r := gin.New()
	r.Use(middlewares.GinLogger(true), corsConfig)

	// repository
	var ()

	// service
	var ()

	// controller
	var ()

	server := &http.Server{
		Addr:    configs.AppConfig.AddressListener(),
		Handler: r,
	}

	go func() {
		logger.Infof("Starting Server on %v", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf(err, "Opening HTTP server: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	logger.Infof("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Shutdown error: %v", err)
	}
}

func init() {
	_, err := configs.Load()
	if err != nil {
		os.Exit(AppExitCode)
	}

	if err != nil {
		fmt.Println(err, "Initializing error messages resource: %v", err)
	}
}
