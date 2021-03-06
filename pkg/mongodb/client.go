package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"time"
)

type Client struct {
	uri      string
	database string
	client   *mongo.Client
}

func NewClient(uri string) (*Client, error) {
	connectionOptions, err := connstring.ParseAndValidate(uri)
	if err != nil {
		return nil, err
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{
		uri:      uri,
		database: connectionOptions.Database,
		client:   client,
	}, nil
}

func (c *Client) Database() *mongo.Database {
	return c.client.Database(c.database)
}

func (c *Client) Collection(collection string) *mongo.Collection {
	return c.Database().Collection(collection)
}
