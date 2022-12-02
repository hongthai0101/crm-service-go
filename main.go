package main

import (
	"context"
	"crm-service-go/config"
	"crm-service-go/datasources"
	"crm-service-go/docs"
	"crm-service-go/pkg/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ctx         context.Context
	mongoClient *mongo.Client
	cancel      context.CancelFunc
)

func init() {
	config.LoadENV()
	utils.InitializeLogger()

	mongoClient, ctx, cancel, _ = datasources.ConnectDB()
}

// @BasePath	/api/v1
func main() {

	utils.Debug(config.GetEnv("LOGNAME", "ssss"))

	defer datasources.Close(mongoClient, ctx, cancel)

	srv := initServer()
	router := srv.Run()

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	_ = router.Run(":" + config.GetEnv("PORT", "3000"))
}
