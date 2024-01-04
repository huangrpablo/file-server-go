package service

import (
	"context"
	"fmt"
	"github.com/file-server-go/storage"
	"github.com/file-server-go/types"
	"github.com/google/uuid"
	"os"
)

type DownloadService struct {
	store  storage.FileStore
	crypto types.Crypto
}

func NewDownloadService(store storage.FileStore, crypto types.Crypto) *DownloadService {
	return &DownloadService{
		store:  store,
		crypto: crypto,
	}
}

func (s *DownloadService) Execute(ctx context.Context, filename string) (filepath string, err error) {
	content, err := s.store.Download(ctx, filename)

	if err != nil {
		return "", fmt.Errorf("failed to download the file: %w", err)
	}

	decrypted, err := s.crypto.Decrypt(content)

	if err != nil {
		return "", fmt.Errorf("failed to decrpyt the file: %w", err)
	}

	// to avoid conflicts from concurrent requests
	filepath = "downloaded/" + filename + uuid.New().String()

	file, err := os.Create(filepath)

	if err != nil {
		return "", fmt.Errorf("failed to create a file: %w", err)
	}

	_, err = file.Write(decrypted)

	if err != nil {
		return "", fmt.Errorf("failed to write the downloaded file: %w", err)
	}

	return filepath, nil
}
