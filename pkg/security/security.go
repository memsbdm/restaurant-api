package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("error during password hash: %w", err)
	}

	return string(bytes), nil
}

func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("error during random string generation: %w", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func SignString(data string, secretKey []byte) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(data))

	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func VerifySignature(data, signature string, secretKey []byte) bool {
	expectedSignature := SignString(data, secretKey)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
