package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"bitbucket-cli/internal/config"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	authHeader string
}

func NewClient(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("missing config")
	}
	auth, err := authHeader(cfg)
	if err != nil {
		return nil, err
	}
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.bitbucket.org/2.0"
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authHeader: auth,
	}, nil
}

func (c *Client) CurrentUser(ctx context.Context) (*User, error) {
	var user User
	if err := c.do(ctx, http.MethodGet, "/user", nil, nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body any, out any) error {
	endpoint := path
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		if query != nil && len(query) > 0 {
			endpoint = c.baseURL + path + "?" + query.Encode()
		} else {
			endpoint = c.baseURL + path
		}
	}

	var payload *bytes.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = bytes.NewReader(data)
	} else {
		payload = bytes.NewReader(nil)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return decodeAPIError(resp)
	}
	defer resp.Body.Close()
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func fetchAll[T any](ctx context.Context, c *Client, path string, query url.Values) ([]T, error) {
	if query == nil {
		query = url.Values{}
	}
	if query.Get("pagelen") == "" {
		query.Set("pagelen", "100")
	}

	var out []T
	nextPath := path
	nextQuery := query
	for {
		var page pagedResponse[T]
		if err := c.do(ctx, http.MethodGet, nextPath, nextQuery, nil, &page); err != nil {
			return nil, err
		}
		out = append(out, page.Values...)
		if page.Next == "" {
			return out, nil
		}
		nextPath = page.Next
		nextQuery = nil
	}
}
