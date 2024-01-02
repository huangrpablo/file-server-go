package types

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"golang.org/x/xerrors"
	"io"
)

type Crypto interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

func NewAES() (Crypto, error) {
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)

	if err != nil {
		return nil, err
	}

	return &AES{key}, nil
}

type AES struct {
	key []byte
}

func (a *AES) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	// Generate a random nonce
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Encrypt the content
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	// Append the nonce to the ciphertext
	ciphertext = append(nonce, ciphertext...)

	return ciphertext, nil
}

func (a *AES) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < 12+block.BlockSize() {
		return nil, xerrors.New("invalid ciphertext")
	}

	// Extract the nonce from the ciphertext
	nonce := ciphertext[:12]
	ciphertext = ciphertext[12:]

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt the content
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
