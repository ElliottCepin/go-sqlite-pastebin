package store

import (
	"context"
)

func Store(ctx context.Context, slug string, content string) error {
	return nil	
}

func Get(ctx context.Context, slug string) (string, error) {
	return "", nil
}

func Delete(ctx context.Context, slug string) error {
	return nil	
}
