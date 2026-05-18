package server

import (
	"context"
)

type Store interface {
	Save(ctx context.Context, slug string, content string) error
	Get(ctx context.Context, slug string) (string, error)
	Delete(ctx context.Context, slug string) error
}
