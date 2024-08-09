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
	"github.com/sirupsen/logrus"
)

var r *router.Router

func init() {
	logrus.Info("Initializing serverless API scraper")

	db := db.NewDynamoDB("sa-east-1")
	productRepo := repository.NewProductRepositoryImpl(db, "products")
	scraper := utils.NewScraperImpl()
	productService := service.NewProductServiceImpl(productRepo, scraper)
	productController := controller.NewProductControllerImpl(productService)
	r = router.NewRouter(productController)
	r.InitRoutes()

	logrus.Info("Serverless API scraper initialized Successfully")
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logrus.Info("Handling request:", req)
	response, err := r.Handler(ctx, req)
	if err != nil {
		logrus.Error("Error handling request:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Internal Server Error"}`,
		}, err
	}
	logrus.Info("Request handled successfully")
	return response, nil
}

func main() {
	logrus.Info("Starting serverless API scraper")
	lambda.Start(Handler)
}
