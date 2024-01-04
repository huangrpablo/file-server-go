package service

import (
	"context"
	"fmt"
	"github.com/file-server-go/storage"
	"github.com/file-server-go/types"
	"os"
)

type UploadService struct {
	store  storage.FileStore
	crypto types.Crypto
}

func NewUploadService(store storage.FileStore, crypto types.Crypto) *UploadService {
	return &UploadService{
		store:  store,
		crypto: crypto,
	}
}

func (s *UploadService) Execute(ctx context.Context, filename string, filepath string) error {
	content, err := os.ReadFile(filepath)

	if err != nil {
		return fmt.Errorf("failed to read the uploaded file: %w", err)
	}

	encrypted, err := s.crypto.Encrypt(content)

	if err != nil {
		return fmt.Errorf("failed to encrypt the file: %w", err)
	}

	err = s.store.Upload(ctx, filename, encrypted)

	if err != nil {
		return fmt.Errorf("failed to put the file to bucket: %w", err)
	}

	return nil
}
