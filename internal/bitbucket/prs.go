package bitbucket

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func (c *Client) ListPullRequests(ctx context.Context, workspace, repo, state string) ([]PullRequest, error) {
	query := url.Values{}
	if state != "" {
		query.Set("state", state)
	}
	query.Set("sort", "-updated_on")
	return fetchAll[PullRequest](ctx, c, fmt.Sprintf("/repositories/%s/%s/pullrequests", workspace, repo), query)
}

func (c *Client) GetPullRequest(ctx context.Context, workspace, repo string, id int) (*PullRequest, error) {
	var out PullRequest
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/repositories/%s/%s/pullrequests/%d", workspace, repo, id), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type CreatePullRequestInput struct {
	Title             string
	Description       string
	SourceBranch      string
	DestinationBranch string
	CloseSourceBranch bool
}

func (c *Client) CreatePullRequest(ctx context.Context, workspace, repo string, in CreatePullRequestInput) (*PullRequest, error) {
	body := map[string]any{
		"title":               in.Title,
		"description":         in.Description,
		"close_source_branch": in.CloseSourceBranch,
		"source": map[string]any{
			"branch": map[string]string{"name": in.SourceBranch},
		},
	}
	if in.DestinationBranch != "" {
		body["destination"] = map[string]any{
			"branch": map[string]string{"name": in.DestinationBranch},
		}
	}

	var out PullRequest
	if err := c.do(ctx, http.MethodPost, fmt.Sprintf("/repositories/%s/%s/pullrequests", workspace, repo), nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func ParsePullRequestID(raw string) (int, error) {
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid pull request id %q", raw)
	}
	return id, nil
}
