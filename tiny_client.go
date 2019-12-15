package mono

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type tinyClient struct {
	token        string
	client       HTTPClient
	unmarshaller Unmarshaller
}

func (c tinyClient) request(ctx context.Context, method, url string, body io.Reader, dst interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	if len(c.token) > 0 {
		req.Header.Add("X-Token", c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %v", err)
	}

	if err := resp.Body.Close(); err != nil {
		return fmt.Errorf("failed to close the body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var derp errorMono
		if err := c.unmarshaller.Unmarshal(bts, &derp); err != nil {
			return fmt.Errorf("failed to unmarshal body: %v", err)
		}

		return fmt.Errorf("mono error: %s", derp.Description)
	}

	if err := c.unmarshaller.Unmarshal(bts, &dst); err != nil {
		return fmt.Errorf("failed to unmarshal body: %v", err)
	}

	return nil
}
