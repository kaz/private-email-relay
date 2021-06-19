package assign

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"golang.org/x/net/publicsuffix"
)

var (
	charset = []byte("abcdefghijklmnopqrstuvwxyz")
)

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Int()%len(charset)]
	}
	return string(result)
}

func effectiveDomain(rawurl string) (string, error) {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	edom, err := publicsuffix.EffectiveTLDPlusOne(parsed.Hostname())
	if err != nil {
		return "", fmt.Errorf("cannot determine effective domain: %w", err)
	}

	return edom, nil
}
