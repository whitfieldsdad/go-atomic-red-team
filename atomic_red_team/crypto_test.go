package atomic_red_team

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAES256GCM(t *testing.T) {
	a := []byte("plaintext")
	password := "password"

	k := KDF([]byte(password), nil)
	c, err := SymmetricEncryptBytes(a, k, AES256GCM)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	b, err := SymmetricDecryptBytes(c, k, AES256GCM)
	assert.Nil(t, err)
	assert.NotNil(t, b)

	// Ensure that the plaintext and decrypted ciphertext are the same.
	assert.Equal(t, a, b)
}

func TestPBKDF2_HMAC_SHA256(t *testing.T) {
	password := "password"
	salt := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	rounds := 100000

	expected := []byte{0xa2, 0x9f, 0xea, 0xf, 0xed, 0x85, 0xc5, 0xb8, 0x61, 0xc, 0x2e, 0x56, 0x97, 0xea, 0x41, 0xb5, 0x58, 0x71, 0x39, 0xe5, 0x8a, 0x38, 0x8e, 0xc, 0x7b, 0x7c, 0xed, 0x30, 0xd4, 0xe6, 0xd8, 0xdf}
	result := PBKDF2_HMAC_SHA256(password, salt, rounds)
	assert.Equal(t, expected, result)
}
