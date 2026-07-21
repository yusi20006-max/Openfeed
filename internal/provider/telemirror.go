package provider

import (
	"context"

	"openfeed/internal/telemirror"
)

type TeleMirror struct {
	client *telemirror.Client
}

func NewTeleMirror() *TeleMirror {

	return &TeleMirror{

		client: telemirror.NewClient(),
	}

}

// LoadChannel fetches the raw HTML widget for a channel through the
// telemirror engine (Google Translate domain-fronting + utls fingerprint).
func (t *TeleMirror) LoadChannel(name string) ([]byte, error) {

	html, err := t.client.FetchHTML(context.Background(), name)
	if err != nil {
		return nil, err
	}

	return []byte(html), nil

}

// FetchDownload proxies an arbitrary media URL (video/audio/document)
// through the same domain-fronting path used for the channel widget,
// sized for large files rather than the small image/thumbnail cap.
func (t *TeleMirror) FetchDownload(rawURL string) ([]byte, string, error) {

	return t.client.FetchDownload(context.Background(), rawURL)

}

// FetchImage proxies a translate.goog-rewritten image URL (avatar or
// post thumbnail) through the same domain-fronting path, capped at the
// smaller image size limit. This lets the browser load images from our
// own server instead of reaching *.translate.goog directly, which is
// blocked on the client's network without a VPN.
func (t *TeleMirror) FetchImage(rawURL string) ([]byte, string, error) {

	return t.client.FetchURL(context.Background(), rawURL)

}
