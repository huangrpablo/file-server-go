package fake

import "github.com/file-server-go/types"

type AES struct{}

func NewAES() types.Crypto {
	return &AES{}
}

func (a *AES) Encrypt(plaintext []byte) ([]byte, error) {
	ciphertext := append(plaintext, 'a')
	return ciphertext, nil
}

func (a *AES) Decrypt(ciphertext []byte) ([]byte, error) {
	return ciphertext[:len(ciphertext)-1], nil
}
