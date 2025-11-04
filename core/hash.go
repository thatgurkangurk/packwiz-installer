package core

import (
	"fmt"
	"strings"

	packwiz "github.com/packwiz/packwiz/core"
)

// PreferredHashList lists hash algorithms in descending order of preference.
var PreferredHashList = []string{
	"murmur2",
	"md5",
	"sha1",
	"sha256",
	"sha512",
}

// MatchHash verifies that the given data matches the provided hash string
// using the specified hashFormat. It supports all hash algorithms known
// to the underlying packwiz library.
func MatchHash(data []byte, hashFormat, expected string) (bool, error) {
	if expected == "" {
		return false, fmt.Errorf("MatchHash: expected hash is empty (format=%s)", hashFormat)
	}

	hasher, err := packwiz.GetHashImpl(hashFormat)
	if err != nil {
		return false, fmt.Errorf("MatchHash: unsupported hash format %q: %w", hashFormat, err)
	}

	if _, err := hasher.Write(data); err != nil {
		return false, fmt.Errorf("MatchHash: failed to hash data (%s): %w", hashFormat, err)
	}

	sum := fmt.Sprintf("%x", hasher.Sum(nil))
	match := strings.EqualFold(expected, sum)

	return match, nil
}
