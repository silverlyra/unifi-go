package api

import (
	"context"
	"fmt"
	"log"
)

// Login holds Unifi OS login credentials for authentication.
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Client) Login(ctx context.Context) error {
	if err := c.getSystem(ctx); err != nil {
		return err
	}

	var self User
	if _, err := c.Do("POST", "/api/auth/login", &c.Credentials, &self); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	c.User = &self
	log.Printf("Logged in as %s (%s)", c.User.Username, c.User.FullName)
	return nil
}
