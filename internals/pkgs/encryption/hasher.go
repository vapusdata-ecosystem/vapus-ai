package encryption

import (
	"crypto/sha256"
	"crypto/sha3"
	"encoding/base64"
)

func GenerateSHA256(data, salt string) string {
	combined := data + salt
	hash := sha256.New()
	hash.Write([]byte(combined))
	hashedBytes := hash.Sum(nil)
	base64Hash := base64.StdEncoding.EncodeToString(hashedBytes)
	return base64Hash
}

// verifyHash compares a given input with a stored hash to see if they match.
func ValidateSHA256(data, salt, storedHash string) bool {
	newHash := GenerateSHA256(data, salt)
	return newHash == storedHash
}

func GenerateSHA3_256(data, salt string) string {
	combined := data + salt
	hash := sha3.New512()
	hash.Write([]byte(combined))
	hashedBytes := hash.Sum(nil)
	base64Hash := base64.StdEncoding.EncodeToString(hashedBytes)
	return base64Hash
}

// verifyHash compares a given input with a stored hash to see if they match.
func ValidateSHA3_256(data, salt, storedHash string) bool {
	newHash := GenerateSHA3_256(data, salt)
	return newHash == storedHash
}

func GenerateSHA3_256FromBytes(data, salt []byte) string {
	combined := string(data) + string(salt)
	hash := sha3.New512()
	hash.Write([]byte(combined))
	hashedBytes := hash.Sum(nil)
	base64Hash := base64.StdEncoding.EncodeToString(hashedBytes)
	return base64Hash
}

// verifyHash compares a given input with a stored hash to see if they match.
func ValidateSHA3_256FromBytes(data, salt []byte, storedHash string) bool {
	newHash := GenerateSHA3_256FromBytes(data, salt)
	return newHash == storedHash
}
