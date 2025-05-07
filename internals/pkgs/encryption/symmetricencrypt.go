package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	"golang.org/x/crypto/pbkdf2"
)

// GenerateSymmKeyWithParams generates a symmetric encryption key with the given parameters.
// It takes a salt string and a variadic list of otherIds as input.
// The function generates a passcode using the otherIds and uses it to derive a key using PBKDF2 algorithm with 4096 iterations and a key length of 32 bytes.
// The derived key is then converted to a string and returned.
func GenerateSymmKeyWithParams(salt string, otherIds ...string) string {
	pass := generatePasscode(otherIds...)
	key := pbkdf2.Key([]byte(pass), []byte(salt), 4096, 32, sha256.New)
	return string(key)
}

// GenerateSymmRandomKey generates a random symmetric encryption key with the given salt.
// It uses the `GenerateRandomString` function from the `dmutils` package to generate a random passcode of length 32.
// The passcode is then used to derive a key using the PBKDF2 algorithm with 4096 iterations and a key length of 32 bytes.
// The derived key is converted to a string and returned.
func GenerateSymmRandomKey(salt string, len int) string {
	pass := dmutils.GenerateRandomString(len)
	key := pbkdf2.Key([]byte(pass), []byte(salt), 4096, len, sha256.New)
	return string(key)
}

func generatePasscode(params ...string) string {
	var pass string
	for _, param := range params {
		pass += param[:len(param)/2]
	}
	return pass
}

func Generate16BitSaltPasscode(params ...string) string {
	var pass string
	for _, param := range params {
		pass += param[:len(param)/2]
	}
	return pass
}

func SymmEncrypt(data, key string) (string, error) {
	// Create a new AES cipher block
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Generate a new GCM (Galois/Counter Mode) cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a random nonce (number used once)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the plaintext using the key, nonce, and GCM mode
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)

	return string(ciphertext), nil
}

func SymmDecrypt(data, key string) (string, error) {
	// Create a new AES cipher block
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Generate a new GCM (Galois/Counter Mode) cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Get the nonce size
	nonceSize := gcm.NonceSize()

	// Extract the nonce from the encrypted data
	nonce, ciphertext := []byte(data[:nonceSize]), []byte(data[nonceSize:])

	// Decrypt the ciphertext using the key, nonce, and GCM mode
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
