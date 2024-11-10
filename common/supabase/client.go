package supabase

import (
	"github.com/supabase-community/supabase-go"
	"sexy_backend/common/log"
)

type Client struct {
	Client *supabase.Client
}

func NewClient(supabaseURL, supabaseKey string) *Client {
	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		log.Error("NewClient err: %v", err)
		panic(err)
	}
	return &Client{
		Client: client,
	}
}

// From returns a database query builder for the specified table
func (c *Client) From(table string) interface{} {
	return c.Client.From(table)
}
