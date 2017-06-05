package weather

import (
	"net/http"
	"time"
)

var (
	defaultRate = time.Second / 60
	defaultBurstLimit = 1
)

type Client struct {
	client *http.Client
	ticker *time.Ticker
	throttle chan time.Time
	rate time.Duration
	burstLimit int
	apiKey string
}

type Option func(c *Client)

func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		rate: defaultRate,
		burstLimit: defaultBurstLimit,
		client: http.DefaultClient,
		apiKey: apiKey,
	}
	for _, opt := range opts {
		opt(c)
	}
	throttle := make(chan time.Time, defaultBurstLimit)
	tick := time.NewTicker(defaultRate)
	go func() {
		for t := range tick.C {
			throttle <- t
		}
		close(throttle)
	}()
	c.ticker = tick
	c.throttle = throttle

	return c
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

func WithRate(rate time.Duration, burstLimit int) Option {
	return func(c *Client) {
		c.rate = rate
		c.burstLimit = burstLimit
	}
}


func (c *Client) Close() {
	c.ticker.Stop()
}