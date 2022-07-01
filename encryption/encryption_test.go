package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryption(t *testing.T) {

	testKey := "iVkmywNjSP0P8K52"

	ce := NewCryptoEngine(testKey)

	encryptedString := ce.Encrypt("string to be encrypted")
	decryptedString := ce.Decrypt(encryptedString)
	assert.Equal(t, "string to be encrypted", decryptedString)
}
