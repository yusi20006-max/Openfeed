package provider

import (
	"fmt"
	"io"
	"net/http"
)

type GoogleTranslate struct {
	client *http.Client
	hosts  []string
}

func NewGoogleTranslate() *GoogleTranslate {

	return &GoogleTranslate{

		client: NewHTTPClient(),

		hosts: Hosts,
	}

}

func (g *GoogleTranslate) LoadChannel(name string) ([]byte, error) {

	var last error

	for _, host := range g.hosts {

		url := fmt.Sprintf(
			"https://%s/s/%s",
			host,
			name,
		)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			last = err
			continue
		}

		req.Header.Set("User-Agent",
			"Mozilla/5.0")

		resp, err := g.client.Do(req)
		if err != nil {
			last = err
			continue
		}

		if resp.StatusCode != 200 {

			resp.Body.Close()
			continue

		}

		body, err := io.ReadAll(resp.Body)

		resp.Body.Close()

		if err != nil {

			last = err
			continue

		}

		return body, nil

	}

	return nil, last

}
