package tests

import (
	"context"
	"github.com/file-server-go/config"
	"github.com/file-server-go/storage/minio"
	"github.com/file-server-go/types"
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

func Test_FileStore_Upload_And_Download(t *testing.T) {
	store := minio.NewFileStore(5242880, false, "taurus-file-store")

	ctx := context.Background()

	file := []byte("abc")
	filename := "test_abc"

	err := store.Upload(ctx, filename, file)
	require.NoError(t, err)

	archive, err := store.Download(ctx, filename)
	require.NoError(t, err)

	require.Equal(t, file, archive)
}

func Test_FileStore_Upload_And_Download_With_Crypto(t *testing.T) {
	aes, err := types.NewAES()
	require.NoError(t, err)

	store := minio.NewFileStore(5242880, false, "taurus-file-store")

	ctx := context.Background()

	file := []byte("abc")
	filename := "test_abc"

	encrypted, err := aes.Encrypt(file)
	require.NoError(t, err)

	err = store.Upload(ctx, filename, encrypted)
	require.NoError(t, err)

	archive, err := store.Download(ctx, filename)
	require.NoError(t, err)

	decrypted, err := aes.Decrypt(archive)
	require.NoError(t, err)

	require.Equal(t, file, decrypted)
}
