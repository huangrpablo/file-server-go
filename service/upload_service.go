package service

import (
	"bytes"
	"fmt"
	"github.com/file-server-go/storage"
	"github.com/file-server-go/types"
	"os"
)

type UploadService struct {
	store  storage.Store
	crypto types.Crypto

	bucketName string
}

func NewUploadService(store storage.Store, crypto types.Crypto, bucketName string) *UploadService {
	return &UploadService{
		store:      store,
		crypto:     crypto,
		bucketName: bucketName,
	}
}

func (s *UploadService) Upload(filename string, filepath string) error {
	content, err := os.ReadFile(filepath)

	if err != nil {
		return fmt.Errorf("failed to read the uploaded file: %w", err)
	}

	encrypted, err := s.crypto.Encrypt(content)

	if err != nil {
		return fmt.Errorf("failed to encrypt the file: %w", err)
	}

	err = s.store.Put(s.bucketName, filename, bytes.NewReader(encrypted))

	if err != nil {
		return fmt.Errorf("failed to put the file to bucket: %w", err)
	}

	return nil
}
