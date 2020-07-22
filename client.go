package namespace

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Client struct {
	c *client.Client
}

func NewClient(c *client.Client) Client {
	return Client{c: c}
}

func (c Client) Namespaces() (personal, other, shared []Namespace, err error) {
	if c.c.State()&imap.AuthenticatedState == 0 {
		return nil, nil, nil, client.ErrNotLoggedIn
	}

	cmd := Command{}
	res := Response{}
	status, err := c.c.Execute(cmd, &res)
	return res.Personal, res.Other, res.Shared, status.Err()
}
