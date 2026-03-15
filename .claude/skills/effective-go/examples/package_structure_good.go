package http

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) Get(url string) ([]byte, error) {
	return nil, nil
}
