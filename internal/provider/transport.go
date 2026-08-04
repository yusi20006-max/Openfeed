package provider

import (
	"net/http"
	"time"
)

func NewHTTPClient() *http.Client {

	return &http.Client{

		Timeout: 20 * time.Second,

		Transport: &http.Transport{

			ForceAttemptHTTP2: true,

			MaxIdleConns: 20,

			MaxIdleConnsPerHost: 10,
		},
	}

}
