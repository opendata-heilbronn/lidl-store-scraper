package main

import (
	"github.com/opendata-heilbronn/lidl-store-scraper/pkg/client"
	"github.com/opendata-heilbronn/lidl-store-scraper/pkg/mongodb"
	"github.com/opendata-heilbronn/lidl-store-scraper/pkg/persistence"
	"github.com/opendata-heilbronn/lidl-store-scraper/pkg/scraper"
	"log"
	"os"
)

type config struct {
	mongodbUrl string
	cronString string
}

func getConfig() config {
	dbUrl := os.Getenv("MONGODB_URL")
	if dbUrl == "" {
		dbUrl = "mongodb://localhost:27017/stores"
	}

	cronString := os.Getenv("SCRAPER_CRON_STRING")
	if cronString == "" {
		cronString = "0 0 20 * * *"
	}

	return config{
		mongodbUrl: dbUrl,
		cronString: cronString,
	}
}

func main() {
	config := getConfig()

	storeClient := client.NewStoreClient("https://meinung.lidl.de/stores_de.json")

	mongoDbClient, err := mongodb.NewClient(config.mongodbUrl)
	if err != nil {
		log.Fatalf("error while connecting to mongodb: %v", err)
	}

	storeRepository, err := persistence.NewStoreRepository(mongoDbClient)
	if err != nil {
		log.Fatalf("error while creating store repository: %v", err)
	}

	scrapeService := scraper.NewScraper(config.cronString, storeClient, storeRepository)
	log.Println("Starting initial scrape")
	scrapeService.Scrape()

	log.Println("Starting scrape cronjob")
	err = scrapeService.StartCronJob()
	if err != nil {
		log.Fatalf("error while starting cronjob: %v", err)
	}
}
