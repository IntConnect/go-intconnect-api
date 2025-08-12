package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-intconnect-api/configs"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/middleware"
	"time"
)

func main() {
	viperConfig := viper.New()
	viperConfig.SetConfigFile(".env")
	viperConfig.AddConfigPath(".")
	viperConfig.AutomaticEnv()
	viperConfig.ReadInConfig()

	// Database Initialization
	databaseCredentials := &configs.DatabaseCredentials{
		DatabaseHost:     viperConfig.GetString("DATABASE_HOST"),
		DatabasePort:     viperConfig.GetString("DATABASE_PORT"),
		DatabaseName:     viperConfig.GetString("DATABASE_NAME"),
		DatabasePassword: viperConfig.GetString("DATABASE_PASSWORD"),
		DatabaseUsername: viperConfig.GetString("DATABASE_USERNAME"),
	}
	fmt.Println(databaseCredentials)
	//databaseInstance := configs.NewDatabaseConnection(databaseCredentials)
	//databaseConnection := databaseInstance.GetDatabaseConnection()

	// Config Initialization
	//validatorInstance, engTranslator := configs.InitializeValidator()

	gin.SetMode(gin.DebugMode)
	// Gin Initialization
	ginEngine := gin.Default()
	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	}))
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(exception.Interceptor())
	rootRouterGroup := ginEngine.Group("/")
	rootRouterGroup.Use(middleware.RequestMetaMiddleware())
	ginError := ginEngine.Run(":8080")
	if ginError != nil {
		panic(ginError)
	}
}
