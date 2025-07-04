package util

import (
	"crypto/rand"
	"strings"

	"github.com/anandvarma/namegen"
)

func GenerateRandomUsername() string {
	ngen := namegen.New()
	return ngen.Get()
}

func GenerateUsername(givenName string, familyName string) string {
	username := givenName + "_" + familyName

	if len(username) < 3 {
		username = GenerateRandomUsername()
	}
	if len(username) > 20 {
		username = username[:20]
	}

	username = strings.ToLower(username)

	username = strings.ReplaceAll(username, " ", "_")
	username = strings.ReplaceAll(username, ".", "")
	username = strings.ReplaceAll(username, "-", "")

	username = strings.ReplaceAll(username, "ö", "o")
	username = strings.ReplaceAll(username, "ü", "u")
	username = strings.ReplaceAll(username, "ç", "c")
	username = strings.ReplaceAll(username, "ş", "s")
	username = strings.ReplaceAll(username, "ğ", "g")
	username = strings.ReplaceAll(username, "ı", "i")

	return username
}

func GenerateRandomPassword() string {
	// Generate a random 32-character password
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, 32)

	// Generate random bytes
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// Fallback to a simple random generation if crypto/rand fails
		for i := range b {
			b[i] = charset[i%len(charset)]
		}
		return string(b)
	}

	// Map random bytes to charset characters
	for i := range b {
		b[i] = charset[randomBytes[i]%byte(len(charset))]
	}

	return string(b)
}
