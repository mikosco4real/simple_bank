package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const alphanumeric = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ&%?#"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min +1)
}

// RandomString generates a random string of length n
func RandomString(n int, s string) string {
	var sb strings.Builder
	k := len(s)

	for i := 0; i < n; i++ {
		c := s[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func GenericRandomString(n int) string {
	return RandomString(n, alphanumeric)
}

// RandomName generates a random (first/last) name
func RandomName() string {
	return strings.Title(RandomString(6, alphabet))
}

// RandomEmail generates a random email
func RandomEmail() string {
	return RandomString(10, alphabet) + "@" + RandomString(5, alphabet) + ".com"
}

// RandomUuid generates a random uuid
func RandomUuid() string {
	return RandomString(16, alphanumeric)
}

func RandomPassword() string {
	return RandomString(10, alphanumeric)
}