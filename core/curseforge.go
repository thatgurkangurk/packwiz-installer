package core

import (
	"context"
	"fmt"
	"os"

	"github.com/carlmjohnson/requests"
)

var (
	cf_api_key         = ""
	cf_api_host        = "https://api.curseforge.com"
	DefaultCurseClient = NewCurseClient(getApiKey())
)

type cfDownloadUrlRes struct {
	Data string `json:"data"`
}

type CurseClient struct {
	apiKey     string
	httpClient *requests.Builder
}

func NewCurseClient(apiKey string) *CurseClient {
	return &CurseClient{
		apiKey: apiKey,
		httpClient: defaultRequestBuilder.
			Clone().
			BaseURL(cf_api_host).
			Header("x-api-key", apiKey),
	}
}

func getApiKey() string {
	key := os.Getenv("CF_API_KEY")
	if key == "" {
		key = cf_api_key
	}
	return key
}

func (c *CurseClient) getJson(ctx context.Context, path string, v any) error {
	if c.apiKey == "" {
		return fmt.Errorf("invalid curseforge api key")
	}

	err := c.httpClient.Path(path).ToJSON(&v).Fetch(context.WithoutCancel(ctx))
	if err != nil {
		return fmt.Errorf("curseforge api: %w", err)
	}
	return nil
}

func (c *CurseClient) GetDownloadUrl(ctx context.Context, d *CurseforgeData) (string, error) {
	path := fmt.Sprintf("/v1/mods/%d/files/%d/download-url", d.ProjectID, d.FileID)
	var resUrl cfDownloadUrlRes
	err := c.getJson(ctx, path, &resUrl)
	if err != nil {
		return "", err
	}
	return resUrl.Data, nil
}
