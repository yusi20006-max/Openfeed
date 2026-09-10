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

// ClientFetchHTMLBefore fetches the channel widget strictly older than
// the supplied Telegram message ID. It is used by the API pagination layer.
func (t *TeleMirror) ClientFetchHTMLBefore(name string, beforeID int) ([]byte, error) {
	html, err := t.client.FetchHTMLBefore(context.Background(), name, beforeID)
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

// FetchImage proxies an avatar or thumbnail through the image-aware
// telemirror path. Image requests use media Accept/Fetch headers and
// reject HTML/error bodies before they reach the browser.
func (t *TeleMirror) FetchImage(rawURL string) ([]byte, string, error) {

	return t.client.FetchImage(context.Background(), rawURL)

}
