package scraper

import (
	"context"
	"fmt"
	"github.com/opendata-heilbronn/lidl-store-scraper/pkg/client"
	"github.com/opendata-heilbronn/lidl-store-scraper/pkg/openinghours"
	"github.com/opendata-heilbronn/lidl-store-scraper/pkg/persistence"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

type Scraper struct {
	cronString  string
	cron        *cron.Cron
	storeClient *client.StoreClient
	storeRepo   *persistence.StoreRepository
}

func NewScraper(
	cronString string,
	storeClient *client.StoreClient,
	storeRepo *persistence.StoreRepository,
) *Scraper {
	return &Scraper{
		cronString:  cronString,
		cron:        cron.New(cron.WithSeconds()),
		storeClient: storeClient,
		storeRepo:   storeRepo,
	}
}

func (s *Scraper) StartCronJob() error {
	_, err := s.cron.AddFunc(s.cronString, s.Scrape)

	if err != nil {
		return err
	}
	s.cron.Run()
	return nil
}

func (s *Scraper) Scrape() {
	scrapeTime := time.Now()
	resp, err := s.storeClient.Scrape()
	if err != nil {
		log.Printf("error while scraping: %v", err)
		return
	}

	storeDocs := s.transformStores(resp, scrapeTime)

	err = s.storeRepo.IngestStores(context.Background(), storeDocs)
	if err != nil {
		log.Printf("error while persisting stores: %v", err)
	}
}

func (s *Scraper) transformStores(in []client.StoreResult, scrapeTime time.Time) []persistence.Store {
	var out []persistence.Store
	for _, inStore := range in {
		openingHours, extraHours, err := openinghours.ParseOpeningHours(inStore.OpeningHours)
		if err != nil {
			log.Printf("error while parsing opening hours for entry %s: %v", inStore.Store, err)
			openingHours = nil
		}

		out = append(out, persistence.Store{
			Store:   inStore.Store,
			Country: inStore.Country,
			ZipCode: fmt.Sprintf("%d", inStore.ZipCode),
			City:    inStore.City,
			Street:  inStore.Street,
			Coordinates: persistence.GeoJsonPoint{
				Type: "Point",
				Coordinates: [2]float64{
					inStore.Longitude, inStore.Latitude,
				},
			},
			ObjectType:   inStore.ObjectType,
			OpeningHours: openingHours,
			ExtraHours:   extraHours,
			ImportTime:   scrapeTime,
		})
	}
	return out
}
