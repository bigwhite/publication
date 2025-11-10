package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserServiceClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

type UserPermissionResponse struct {
	CanCreate bool `json:"can_create"`
}

func (c *UserServiceClient) CanCreateLink(ctx context.Context, userID string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/users/%s/permissions", c.BaseURL, userID), nil)
	if err != nil {
		return false, err
	}
	// 在真实应用中，这里应该有认证头

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("user service returned status %d", resp.StatusCode)
	}

	var permResp UserPermissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&permResp); err != nil {
		return false, err
	}

	return permResp.CanCreate, nil
}
