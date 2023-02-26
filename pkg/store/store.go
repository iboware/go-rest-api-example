package store

import (
	"context"

	"github.com/iboware/go-rest-api-example/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen --build_flags=--mod=mod -destination=./mocks/mock_store.go -package=mocks . Store
type Store interface {
	Find(ctx context.Context, filter primitive.D) ([]model.Record, error)
}
