package bitbucket

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) ListRepositories(ctx context.Context, workspace string) ([]Repository, error) {
	query := url.Values{}
	query.Set("sort", "full_name")
	return fetchAll[Repository](ctx, c, fmt.Sprintf("/repositories/%s", workspace), query)
}

func (c *Client) GetRepository(ctx context.Context, workspace, repo string) (*Repository, error) {
	var out Repository
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/repositories/%s/%s", workspace, repo), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateRepository(ctx context.Context, workspace, repo string, isPrivate bool) (*Repository, error) {
	var out Repository
	body := map[string]any{
		"scm":        "git",
		"is_private": isPrivate,
	}
	if err := c.do(ctx, http.MethodPost, fmt.Sprintf("/repositories/%s/%s", workspace, repo), nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
