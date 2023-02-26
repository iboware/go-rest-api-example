package store

import (
	"context"

	"github.com/iboware/go-rest-api-example/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStore struct {
	client   *mongo.Client
	database string
	table    string
}

func NewMongoStore(client *mongo.Client, database, table string) *MongoStore {
	return &MongoStore{
		client:   client,
		database: database,
		table:    table,
	}
}

func (s *MongoStore) Find(ctx context.Context, filter primitive.D) ([]model.Record, error) {

	// find the matching documents from the collection using filter
	cursor, err := s.client.Database(s.database).Collection(s.table).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	results := make([]model.Record, 0)
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
