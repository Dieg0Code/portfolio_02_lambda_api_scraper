package router

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/dieg0code/serverles-api-scraper/api/controller"
	"github.com/gin-gonic/gin"
)

type Router struct {
	ProductController controller.ProductController
	ginLambda         *ginadapter.GinLambda
}

func NewRouter(productController controller.ProductController) *Router {
	return &Router{
		ProductController: productController,
	}
}

func (r *Router) InitRoutes() *Router {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to serverless API scraper",
		})
	})

	baseRoute := router.Group("/api/v1")
	{
		productRoute := baseRoute.Group("/products")
		{
			productRoute.GET("", r.ProductController.GetAll)
			productRoute.GET("/:productId", r.ProductController.GetByID)
			productRoute.POST("", r.ProductController.UpdateData)
		}
	}

	r.ginLambda = ginadapter.New(router)
	return r
}

func (r *Router) Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return r.ginLambda.ProxyWithContext(ctx, req)
}
