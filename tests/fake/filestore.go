package fake

import (
	"context"
	"github.com/file-server-go/storage"
	"github.com/file-server-go/types"
)

type FileStore map[string][]byte

func NewFileStore() storage.FileStore {
	return &FileStore{}
}

func (s *FileStore) Upload(ctx context.Context, filename string, file []byte) error {
	(*s)[filename] = file
	return nil
}

func (s *FileStore) Download(ctx context.Context, filename string) ([]byte, error) {
	f, ok := (*s)[filename]

	if !ok {
		return nil, types.ErrFileNotFound
	}

	return f, nil
}

func GetFile(store storage.FileStore, filename string) ([]byte, bool) {
	filestore, ok := store.(*FileStore)
	if !ok {
		return nil, false
	}

	file, ok := (*filestore)[filename]
	return file, ok
}
