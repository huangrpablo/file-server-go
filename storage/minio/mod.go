package minio

import (
	"context"
	"fmt"
	"github.com/file-server-go/storage"
	"github.com/file-server-go/types"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

func New(endpoint, accessKey, secretKey string, chunkSize uint64) (storage.Store, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a minio instance: %w", err)
	}

	return &Store{
		Client:    client,
		chunkSize: chunkSize,
	}, nil
}

type Store struct {
	*minio.Client
	chunkSize uint64 // configured to upload files in chunks
}

func (s *Store) Put(bucket, key string, file io.Reader) error {
	options := minio.PutObjectOptions{PartSize: s.chunkSize}

	_, err := s.PutObject(context.Background(), bucket, key, file, -1, options)

	if err != nil {
		return fmt.Errorf("failed to put object(%v) in bucket(%v): %w", key, bucket, err)
	}

	return nil
}

func (s *Store) Get(bucket, key string) (io.Reader, error) {
	file, err := s.GetObject(context.Background(), bucket, key, minio.GetObjectOptions{})

	if err != nil {
		errResp := minio.ToErrorResponse(err)

		if errResp.Code == "NoSuchKey" {
			return nil, fmt.Errorf("failed to get object(%v) from bucket(%v): %w", key, bucket, types.ErrFileNotFound)
		}

		return nil, fmt.Errorf("failed to get object(%v) from bucket(%v): %w", key, bucket, err)
	}

	return file, nil
}
