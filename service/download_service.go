package service

import (
	"fmt"
	"github.com/file-server-go/storage"
	"github.com/file-server-go/types"
	"io"
	"os"
)

type DownloadService struct {
	store  storage.Store
	crypto types.Crypto

	bucketName string
}

func NewDownloadService(store storage.Store, crypto types.Crypto, bucketName string) *DownloadService {
	return &DownloadService{
		store:      store,
		crypto:     crypto,
		bucketName: bucketName,
	}
}

func (s *DownloadService) Download(filename string) (filepath string, err error) {
	archive, err := s.store.Get(s.bucketName, filename)

	if err != nil {
		return "", fmt.Errorf("failed to get the file from bucket: %w", err)
	}

	filepath = "downloaded/" + filename
	file, err := os.Create(filepath)

	if err != nil {
		return "", fmt.Errorf("failed to create a file: %w", err)
	}

	_, err = io.Copy(file, archive)

	if err != nil {
		return "", fmt.Errorf("failed to write the downloaded file: %w", err)
	}

	return filepath, nil
}
