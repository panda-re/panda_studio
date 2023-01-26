package repos

import (
	"context"

	"github.com/panda-re/panda_studio/internal/db"
)

type Repository [T any] interface {
	FindAll(ctx context.Context) ([]T, error)
	FindOne(ctx context.Context, id db.ObjectID) (*T, error)
}