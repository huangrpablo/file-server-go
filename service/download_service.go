package service

import (
	"context"
	"fmt"
	"github.com/file-server-go/storage"
	"github.com/file-server-go/types"
	"github.com/google/uuid"
	"os"
	"path/filepath"
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

func (s *DownloadService) Execute(ctx context.Context, filename string) (path string, err error) {
	content, err := s.store.Download(ctx, filename)

	if err != nil {
		return "", fmt.Errorf("failed to download the file: %w", err)
	}

	decrypted, err := s.crypto.Decrypt(content)

	if err != nil {
		return "", fmt.Errorf("failed to decrpyt the file: %w", err)
	}

	return s.saveFile(filename, decrypted)
}

func (s *DownloadService) saveFile(filename string, content []byte) (path string, err error) {
	// to avoid conflicts from concurrent requests
	path = "downloaded/" + filename + uuid.New().String()

	err = os.MkdirAll(filepath.Dir(path), 0750)

	if err != nil {
		return "", fmt.Errorf("failed to create the dir: %w", err)
	}

	file, err := os.Create(path)

	if err != nil {
		return "", fmt.Errorf("failed to create a file: %w", err)
	}

	_, err = file.Write(content)

	if err != nil {
		return "", fmt.Errorf("failed to write the downloaded file: %w", err)
	}

	return path, nil
}
