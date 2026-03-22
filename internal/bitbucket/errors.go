package bitbucket

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type apiErrorBody struct {
	Type  string `json:"type"`
	Error struct {
		Message string `json:"message"`
		Detail  string `json:"detail"`
	} `json:"error"`
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("bitbucket api returned status %d", e.StatusCode)
	}
	return fmt.Sprintf("bitbucket api returned status %d: %s", e.StatusCode, e.Message)
}

func decodeAPIError(resp *http.Response) error {
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	msg := strings.TrimSpace(string(body))

	var parsed apiErrorBody
	if err := json.Unmarshal(body, &parsed); err == nil {
		switch {
		case parsed.Error.Message != "":
			msg = parsed.Error.Message
		case parsed.Error.Detail != "":
			msg = parsed.Error.Detail
		}
	}

	return &APIError{StatusCode: resp.StatusCode, Message: msg}
}
