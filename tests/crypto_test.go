package tests

import (
	"github.com/file-server-go/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AES(t *testing.T) {
	aes, err := types.NewAES()
	require.NoError(t, err)

	msg := []byte("abc")

	encrypted, err := aes.Encrypt(msg)
	require.NoError(t, err)

	decrypted, err := aes.Decrypt(encrypted)
	require.NoError(t, err)

	require.Equal(t, msg, decrypted)
}
