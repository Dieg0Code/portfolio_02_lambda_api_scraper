package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dieg0code/scraper/src/repository"
	"github.com/dieg0code/scraper/src/scraper"
	"github.com/dieg0code/scraper/src/service"
	"github.com/dieg0code/shared/db"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

var scraperService service.ScraperService

func init() {
	logrus.Info("Initializing scraper")

	region := "sa-east-1"
	tableName := "Products"

	db := db.NewDynamoDB(region)

	scraperRepo := repository.NewScraperRepositoryImpl(db, tableName)

	collector := colly.NewCollector()
	scraper := scraper.NewScraperImpl(collector)

	scraperService = service.NewScraperServiceImpl(scraper, scraperRepo)
}

func handleRequest(ctx context.Context) (bool, error) {
	logrus.Info("Handling request")
	success, err := scraperService.GetProducts()
	if err != nil {
		logrus.WithError(err).Error("Error handling request")
		return false, err
	}

	logrus.Info("Request handled successfully")
	return success, nil
}

func main() {
	lambda.Start(handleRequest)
}
