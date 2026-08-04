package provider

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Direct struct {
	client *http.Client
}

func NewDirect() *Direct {

	return &Direct{
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}

}

func (d *Direct) LoadChannel(name string) ([]byte, error) {

	url := fmt.Sprintf(
		"https://t.me/s/%s",
		name,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent",
		"Mozilla/5.0")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)

}
