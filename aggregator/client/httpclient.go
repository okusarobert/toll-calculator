package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/okusarobert/toll-calculator/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) Client {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	httpc := http.DefaultClient
	endpoint := fmt.Sprintf("%s/invoice?obu=%d", c.Endpoint, id)
	fmt.Println("hitting ...", endpoint)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the service responded with a none 200 status code %d", resp.StatusCode)
	}

	var inv types.Invoice

	if err := json.NewDecoder(resp.Body).Decode(&inv); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &inv, nil
}

func (c *HTTPClient) Aggregate(ctx context.Context, payload *types.AggregateRequest) error {
	httpc := http.DefaultClient
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint+"/aggregate", bytes.NewReader(b))
	if err != nil {
		return err
	}
	resp, err := httpc.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the service responded with a none 200 status code %d", resp.StatusCode)
	}
	resp.Body.Close()
	return nil
}
