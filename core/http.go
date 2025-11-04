package core

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/carlmjohnson/requests"
)

const userAgent = "packwiz-installer by thatgurkangurk (https://github.com/thatgurkangurk/packwiz-installer)"

var defaultRequestBuilder = newRequestBuilder(nil)

// newRequestBuilder returns a requests.Builder configured with a client and
// the package user-agent. If a nil client is passed, http.DefaultClient is used.
func newRequestBuilder(c *http.Client) *requests.Builder {
	if c == nil {
		c = http.DefaultClient
	}
	return requests.New().
		Client(c).
		UserAgent(userAgent)
}

// httpGetJson fetches url and decodes the JSON/TOML/whatever body into v.
// v must be a pointer to the destination value (e.g. *T or *map[string]any).
// The caller is responsible for providing a context with a timeout/cancel if desired.
func httpGetJson(ctx context.Context, c *http.Client, url string, v any) error {
	if v == nil {
		return fmt.Errorf("httpGetJson: destination v must be a non-nil pointer (url=%s)", url)
	}

	// ensure we have a usable client
	builder := defaultRequestBuilder.Clone().Client(c)

	// requests.ToJSON expects the destination (pointer) directly, not &v
	if err := builder.BaseURL(url).ToJSON(v).Fetch(ctx); err != nil {
		return fmt.Errorf("httpGetJson: failed fetching %s: %w", url, err)
	}
	return nil
}

// httpGetBytes fetches url and returns the response body as bytes.
// The caller can control cancellation/timeout using ctx.
func httpGetBytes(ctx context.Context, c *http.Client, url string) ([]byte, error) {
	buf := &bytes.Buffer{}

	builder := defaultRequestBuilder.Clone().Client(c)

	if err := builder.BaseURL(url).ToBytesBuffer(buf).Fetch(ctx); err != nil {
		return nil, fmt.Errorf("httpGetBytes: failed fetching %s: %w", url, err)
	}
	return buf.Bytes(), nil
}

// httpGetValidBytes fetches url, verifies its hash against hashFormat/hash,
// and returns the bytes if the hash matches.
func httpGetValidBytes(ctx context.Context, c *http.Client, url string, hashFormat string, hash string) ([]byte, error) {
	data, err := httpGetBytes(ctx, c, url)
	if err != nil {
		return nil, err
	}

	valid, err := MatchHash(data, hashFormat, hash)
	if err != nil {
		return nil, fmt.Errorf("httpGetValidBytes: validating %s: %w", url, err)
	}
	if !valid {
		return nil, fmt.Errorf("httpGetValidBytes: download hash mismatch for %s", url)
	}
	return data, nil
}
