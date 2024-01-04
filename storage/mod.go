package storage

import (
	"context"
)

type FileStore interface {
	Upload(ctx context.Context, filename string, file []byte) error
	Download(ctx context.Context, filename string) ([]byte, error)
}
