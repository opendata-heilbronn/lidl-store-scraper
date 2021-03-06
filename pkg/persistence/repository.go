package persistence

import (
	"context"
	"github.com/opendata-heilbronn/lidl-store-scraper/pkg/mongodb"
	"github.com/opendata-heilbronn/lidl-store-scraper/pkg/openinghours"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type StoreRepository struct {
	client *mongodb.Client
}

type Store struct {
	Id           primitive.ObjectID             `bson:"_id"`
	Store        string                         `bson:"store"`
	Country      string                         `bson:"country"`
	ZipCode      string                         `bson:"zip"`
	City         string                         `bson:"city"`
	Street       string                         `bson:"street"`
	Coordinates  GeoJsonPoint                   `bson:"coordinates"`
	ObjectType   string                         `bson:"object_type"`
	OpeningHours openinghours.OpeningHours      `bson:"opening_hours"`
	ExtraHours   openinghours.ExtraOpeningHours `bson:"extra_hours"`
	ImportTime   time.Time                      `bson:"import_time"`
}

type GeoJsonPoint struct {
	// Always "Point"
	Type        string     `bson:"type"`
	Coordinates [2]float64 `bson:"coordinates"`
}

func NewStoreRepository(client *mongodb.Client) (*StoreRepository, error) {
	return &StoreRepository{
		client: client,
	}, nil
}

func (r *StoreRepository) IngestStores(ctx context.Context, stores []Store) error {
	storeDocuments := make([]interface{}, 0, len(stores))
	for _, store := range stores {
		store.Id = primitive.NewObjectID()
		storeDocuments = append(storeDocuments, store)
	}
	_, err := r.client.Collection("stores").InsertMany(ctx, storeDocuments, options.InsertMany().SetOrdered(false))
	if err != nil {
		return err
	}
	return nil
}
