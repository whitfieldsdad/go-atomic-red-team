package atomic_red_team

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"
)

const (
	AES256GCM                           = "aes-256-gcm"
	DefaultSymmetricEncryptionAlgorithm = AES256GCM
)

func SymmetricEncryptBytes(plaintext []byte, password, algorithm string) ([]byte, error) {
	if algorithm == "" {
		algorithm = DefaultSymmetricEncryptionAlgorithm
	}
	if algorithm == AES256GCM {
		return AES256GCMEncryptBytes(plaintext, password)
	} else {
		return nil, errors.New("unsupported algorithm")
	}
}

func SymmetricDecryptBytes(ciphertext []byte, password, algorithm string) ([]byte, error) {
	if algorithm == "" {
		algorithm = DefaultSymmetricEncryptionAlgorithm
	}
	if algorithm == AES256GCM {
		return AES256GCMDecryptBytes(ciphertext, password)
	} else {
		return nil, errors.New("unsupported algorithm")
	}
}

func AES256GCMEncryptBytes(plaintext []byte, password string) ([]byte, error) {
	key := []byte(password)
	aes, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cipher")
	}
	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create gcm")
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.Wrap(err, "failed to read random bytes")
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func AES256GCMDecryptBytes(ciphertext []byte, password string) ([]byte, error) {
	aes, err := aes.NewCipher([]byte(password))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cipher")
	}
	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create gcm")
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt")
	}
	return plaintext, nil
}

const (
	PBKDF2Rounds     = 100000
	PBKDF2           = "pbkdf2"
	PBKDF2SaltLength = 8
)

func PBKDF2_HMAC_SHA256(password string, salt []byte, rounds int) []byte {
	if salt == nil {
		salt = GenerateSalt(PBKDF2SaltLength)
	}
	return pbkdf2.Key([]byte(password), salt, rounds, 32, sha256.New)
}

func KDF(password, salt []byte) []byte {
	return PBKDF2_HMAC_SHA256(string(password), salt, PBKDF2Rounds)
}

func GenerateSalt(bytes int) []byte {
	salt := make([]byte, bytes)
	rand.Read(salt)
	return salt
}
