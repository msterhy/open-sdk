package esx

import "github.com/elastic/go-elasticsearch/v8"

type Client struct {
	cfg         elasticsearch.Config
	TypedClient *elasticsearch.TypedClient
	Client      *elasticsearch.Client
}

type Option func(*Client)

var (
	client *Client
)

func NewClient(opts ...Option) (*Client, error) {
	client = &Client{}
	for _, opt := range opts {
		opt(client)
	}
	newClient, err := elasticsearch.NewTypedClient(client.cfg)
	if err != nil {
		return nil, err
	}
	client.TypedClient = newClient
	client.Client, err = elasticsearch.NewClient(client.cfg)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func WithConfig(cfg elasticsearch.Config) Option {
	return func(c *Client) {
		c.cfg = cfg
	}
}

func WithAddress(address []string) Option {
	return func(c *Client) {
		c.cfg.Addresses = address
	}
}

func WithUsername(username string) Option {
	return func(c *Client) {
		c.cfg.Username = username
	}
}

func WithPassword(password string) Option {
	return func(c *Client) {
		c.cfg.Password = password
	}
}

func WithAPIKey(apiKey string) Option {
	return func(c *Client) {
		c.cfg.APIKey = apiKey
	}
}

func GetClient() *Client {
	return client
}
