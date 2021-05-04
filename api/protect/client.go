package protect

import (
	"context"
	"fmt"

	"github.com/silverlyra/unifi-go/api"
)

type Client struct {
	*api.Client

	Bootstrap *Bootstrap
}

func New(base *api.Client) (*Client, error) {
	client := &Client{
		Client: base,
	}

	return client, nil
}

func (c *Client) Load(ctx context.Context) error {
	if c.User == nil {
		if err := c.Login(ctx); err != nil {
			return err
		}
	}

	var data Bootstrap
	if _, err := c.Do("GET", "/proxy/protect/api/bootstrap", nil, &data); err != nil {
		return fmt.Errorf("failed to get Protect bootstrap data: %w", err)
	}

	c.Bootstrap = &data
	return nil
}
