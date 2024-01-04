package tests

import (
	"context"
	"github.com/file-server-go/service"
	"github.com/file-server-go/tests/fake"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_Service_Upload(t *testing.T) {
	aes := fake.NewAES()
	store := fake.NewFileStore()
	ctx := context.Background()

	// upload
	upload := service.NewUploadService(store, aes)

	file := []byte("abc")
	filename := "test_abc"
	filepath := "files/test_abc"

	err := upload.Execute(ctx, filename, filepath)
	require.NoError(t, err)

	content, ok := fake.GetFile(store, filename)
	require.Equal(t, true, ok)

	encrypted, err := aes.Encrypt(file)
	require.NoError(t, err)
	require.Equal(t, encrypted, content)

	// download
	download := service.NewDownloadService(store, aes)

	filepath, err = download.Execute(ctx, filename)
	require.NoError(t, err)

	content, err = os.ReadFile(filepath)
	require.NoError(t, err)

	require.Equal(t, file, content)
}
