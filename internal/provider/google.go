package provider

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 20 * time.Second,
}

func LoadChannel(name string) ([]byte, error) {

	url := fmt.Sprintf("https://t.me/s/%s", name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent",
		"Mozilla/5.0 (Android 12) AppleWebKit/537.36 Chrome/138 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
