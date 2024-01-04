package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/file-server-go/storage"
	"github.com/file-server-go/types"
	"github.com/minio/minio-go/v7"
	"io"
)

func NewFileStore(chunkSize uint64, uploadInChunk bool, bucketName string) storage.FileStore {
	return &FileStore{
		chunkSize:     chunkSize,
		uploadInChunk: uploadInChunk,
		bucketName:    bucketName,
	}
}

type FileStore struct {
	chunkSize     uint64 // configured to upload files in chunks
	uploadInChunk bool
	bucketName    string
}

func (s *FileStore) Upload(ctx context.Context, filename string, file []byte) error {
	options := minio.PutObjectOptions{PartSize: s.chunkSize, DisableMultipart: !s.uploadInChunk}

	_, err := client.PutObject(ctx, s.bucketName, filename, bytes.NewReader(file), -1, options)

	if err != nil {
		return fmt.Errorf("failed to put object(%v) in bucket(%v): %w", filename, s.bucketName, err)
	}

	return nil
}

func (s *FileStore) Download(ctx context.Context, filename string) ([]byte, error) {
	file, err := s.download(ctx, filename)

	if err != nil {
		errResp := minio.ToErrorResponse(err)

		if errResp.Code == "NoSuchKey" {
			return nil, fmt.Errorf("failed to get object(%v) from bucket(%v): %w", filename, s.bucketName, types.ErrFileNotFound)
		}

		return nil, fmt.Errorf("failed to get object(%v) from bucket(%v): %w", filename, s.bucketName, err)
	}

	return file, nil
}

func (s *FileStore) download(ctx context.Context, filename string) ([]byte, error) {
	archive, err := client.GetObject(ctx, s.bucketName, filename, minio.GetObjectOptions{})

	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(archive)

	if err != nil {
		return nil, err
	}

	return content, nil
}
