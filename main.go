package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dieg0code/serverles-api-scraper/api/controller"
	"github.com/dieg0code/serverles-api-scraper/api/db"
	"github.com/dieg0code/serverles-api-scraper/api/repository"
	"github.com/dieg0code/serverles-api-scraper/api/router"
	"github.com/dieg0code/serverles-api-scraper/api/service"
	"github.com/dieg0code/serverles-api-scraper/api/utils"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var r *router.Router

func init() {
	// load env vars
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Error("Error loading .env file")
	}

	db := db.NewDynamoDB("sa-east-1")
	productRepo := repository.NewProductRepositoryImpl(db, "products")
	scraper := utils.NewScraperImpl()
	productService := service.NewProductServiceImpl(productRepo, scraper)
	productController := controller.NewProductControllerImpl(productService)
	r = router.NewRouter(productController)
	r.InitRoutes()
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return r.Handler(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
