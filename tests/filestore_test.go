package tests

import (
	"context"
	"github.com/file-server-go/config"
	"github.com/file-server-go/storage/minio"
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	conf, err := config.Load("../config/config.yaml")
	if err != nil {
		panic(err)
	}

	err = minio.Init(conf.Minio.Endpoint, conf.Minio.AccessKey, conf.Minio.SecretKey)
	if err != nil {
		panic(err)
	}
}

func Test_Upload_And_Download(t *testing.T) {
	store := minio.NewFileStore(5242880, "taurus-file-store")

	ctx := context.Background()

	file := []byte("abc")
	filename := "test_abc"

	err := store.Upload(ctx, filename, file)
	require.NoError(t, err)

	archive, err := store.Download(ctx, filename)
	require.NoError(t, err)

	require.Equal(t, filename, archive)
}
