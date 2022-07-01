package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

type CryptoEngine struct {
	key string
}

func NewCryptoEngine(key string) CryptoEngine {
	return CryptoEngine{
		key: key,
	}
}

func (ce *CryptoEngine) Encrypt(content string) string {

	plaintext := []byte(content)

	block, err := aes.NewCipher([]byte(ce.key))
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

func (ce *CryptoEngine) Decrypt(encryptedText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(encryptedText)

	block, err := aes.NewCipher([]byte(ce.key))
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext)
}
