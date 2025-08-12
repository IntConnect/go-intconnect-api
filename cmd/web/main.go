package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-intconnect-api/configs"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/middleware"
	"go.uber.org/fx"
	"time"
)

// --- Provider untuk Viper config ---
func NewViperConfig() *viper.Viper {
	viperConfig := viper.New()
	viperConfig.SetConfigFile(".env")
	viperConfig.AddConfigPath(".")
	viperConfig.AutomaticEnv()
	if err := viperConfig.ReadInConfig(); err != nil {
		panic(err)
	}
	return viperConfig
}

// --- Provider untuk Database Credentials ---
func NewDatabaseCredentials(viperConfig *viper.Viper) *configs.DatabaseCredentials {
	return &configs.DatabaseCredentials{
		DatabaseHost:     viperConfig.GetString("DATABASE_HOST"),
		DatabasePort:     viperConfig.GetString("DATABASE_PORT"),
		DatabaseName:     viperConfig.GetString("DATABASE_NAME"),
		DatabasePassword: viperConfig.GetString("DATABASE_PASSWORD"),
		DatabaseUsername: viperConfig.GetString("DATABASE_USERNAME"),
	}
}

// --- Provider untuk Gin Engine ---
func NewGinEngine() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	ginEngine := gin.Default()
	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(exception.Interceptor())
	ginEngineRoot := ginEngine.Group("/")
	ginEngineRoot.Use(middleware.RequestMetaMiddleware())

	return ginEngine
}

// --- Invoker ---
func Run(fxLifecycle fx.Lifecycle, ginEngine *gin.Engine) {
	fxLifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Starting server...")
			go func() {
				if err := ginEngine.Run(":8080"); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping server...")
			return nil
		},
	})
}

func main() {
	fxContainer := fx.New(
		// Provider
		fx.Provide(
			NewViperConfig,
			NewDatabaseCredentials,
			// configs.NewDatabaseConnection, // kalau mau tambahkan DB
			NewGinEngine,
		),
		// Invoker
		fx.Invoke(Run),
	)

	fxContainer.Run()
}
