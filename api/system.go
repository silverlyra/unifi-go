package api

import (
	"context"
	"fmt"
)

// System describes the Unifi OS installation responding to API requests.
type System struct {
	Name     string   `json:"name"`
	Hardware Hardware `json:"hardware"`
}

// Hardware describes the Unifi device responding to API requests.
type Hardware struct {
	Type HardwareType `json:"shortname"`
}

// HardwareType is the Unifi "shortname" describing the device model.
type HardwareType string

const (
	UnifiDreamMachinePro HardwareType = "UDMPRO"
)

func (c *Client) getSystem(ctx context.Context) error {
	var sys System
	res, err := c.Do("GET", "/api/system", nil, &sys)

	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	token := res.Header.Get("X-Csrf-Token")
	if token == "" {
		return fmt.Errorf("no X-CSRF-Token header in response")
	}
	c.csrfToken = token

	c.System = &sys
	return nil
}
