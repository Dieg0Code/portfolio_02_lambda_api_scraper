package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dieg0code/api-users/controllers"
	"github.com/dieg0code/api-users/repository"
	"github.com/dieg0code/api-users/router"
	"github.com/dieg0code/api-users/services"
	"github.com/dieg0code/api-users/utils"
	"github.com/dieg0code/shared/db"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

var r *router.Router

func init() {
	logrus.Info("Initializing serverless API users")

	region := "sa-east-1"
	tableName := "Users"

	// Instance Database
	db := db.NewDynamoDB(region)

	// Instance repository
	userRepo := repository.NewUserRepositoryImpl(db, tableName)

	validator := validator.New()
	passwordHaher := utils.NewPasswordHasher()
	jwtUtils := utils.NewJWTUtils()

	// Instance Service
	userService := services.NewUserServiceImpl(userRepo, validator, passwordHaher, jwtUtils)

	// Instance controller
	userController := controllers.NewUserControllerImpl(userService)

	r = router.NewRouter(userController)
	r.InitRoutes()

	logrus.Info("Serverless API users initialized Successfully")
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logrus.Info("Handling request:", req.Path)
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
	logrus.Info("Starting serverless API users")
	lambda.Start(Handler)
}
