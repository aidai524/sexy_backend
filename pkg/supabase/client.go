package supabase

import (
	"github.com/supabase-community/supabase-go"
)

type Client struct {
	Client *supabase.Client
}

func NewClient(supabaseURL, supabaseKey string) (*Client, error) {
	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client: client,
	}, nil
}

// From returns a database query builder for the specified table
func (c *Client) From(table string) interface{} {
	return c.Client.From(table)
}
