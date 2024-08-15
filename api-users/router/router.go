package router

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/dieg0code/api-users/controllers"
	"github.com/gin-gonic/gin"
)

type Router struct {
	UserController controllers.UserController
	ginLambda      *ginadapter.GinLambda
}

func NewRouter(userController controllers.UserController) *Router {
	return &Router{
		UserController: userController,
	}
}

func (r *Router) InitRoutes() *Router {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to serverless API users",
		})
	})

	baseRoute := router.Group("/api/v1")
	{
		userRoute := baseRoute.Group("/users")
		{
			userRoute.GET("", r.UserController.GetAllUsers)
			userRoute.GET("/:userID", r.UserController.GetUserByID)
			userRoute.POST("", r.UserController.RegisterUser)
		}
	}

	r.ginLambda = ginadapter.New(router)
	return r
}

func (r *Router) Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return r.ginLambda.ProxyWithContext(ctx, req)
}
