package telemirror

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// FetchImage proxies an image URL using image-specific request headers.
// The channel HTML path uses document headers, but Telegram thumbnails and
// avatars are media resources; asking the Translate edge for a document can
// return an HTML wrapper/error instead of the image bytes.
func (c *Client) FetchImage(ctx context.Context, rawURL string) ([]byte, string, error) {
	u, err := validateSafeURL(ctx, rawURL)
	if err != nil {
		return nil, "", err
	}

	frontHost := toTranslateGoogHost(u.Host)
	buildFrontURL := func(sl, tl string) string {
		fu := *u
		fu.Host = frontHost
		q := fu.Query()
		q.Set("_x_tr_sl", sl)
		q.Set("_x_tr_tl", tl)
		q.Set("_x_tr_hl", "en")
		q.Set("_x_tr_pto", "wapp")
		fu.RawQuery = q.Encode()
		return fu.String()
	}

	var lastErr error
	for i, ap := range c.orderedAttempts() {
		if i > 0 {
			select {
			case <-ctx.Done():
				return nil, "", ctx.Err()
			case <-time.After(200 * time.Millisecond):
			}
		}

		frontURL := buildFrontURL(ap.sl, ap.tl)
		body, ctype, status, err := c.fetchImageOnce(ctx, ap, frontURL, frontHost)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d (%s): %w", i+1, ap.label(), err)
			continue
		}
		if status == http.StatusOK {
			c.markSuccess(attemptIndex(ap))
			return body, ctype, nil
		}
		lastErr = fmt.Errorf("attempt %d (%s) status %d", i+1, ap.label(), status)
	}

	// If the URL is already a usable translate.goog URL, or the generated
	// proxy route is unavailable, retain the safe direct fallback.
	if frontHost != u.Host {
		body, ctype, status, ferr := c.fetchImageOnce(ctx,
			proxyAttempt{sni: sniUseHost, fp: fingerprints[1]}, u.String(), u.Host)
		if ferr == nil && status == http.StatusOK {
			return body, ctype, nil
		}
		if ferr != nil {
			lastErr = fmt.Errorf("direct fallback: %w", ferr)
		} else {
			lastErr = fmt.Errorf("direct fallback status %d", status)
		}
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("telemirror: all image attempts exhausted")
	}
	return nil, "", lastErr
}

func setImageHeaders(req *http.Request, ua string, fronted bool) {
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,fa;q=0.8")
	req.Header.Set("Sec-Fetch-Dest", "image")
	req.Header.Set("Sec-Fetch-Mode", "no-cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	if fronted {
		req.Header.Set("Referer", "https://translate.google.com/")
	}
}

func (c *Client) fetchImageOnce(ctx context.Context, ap proxyAttempt, rawURL, hostHeader string) ([]byte, string, int, error) {
	transport := transportFor(ap, hostHeader)
	defer transport.CloseIdleConnections()
	httpClient := &http.Client{Transport: transport, Timeout: requestTimeout, CheckRedirect: safeCheckRedirect}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", 0, err
	}
	req.Host = hostHeader
	ua := userAgents[0]
	setImageHeaders(req, ua, ap.ip != "")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, "", 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize+1))
	if err != nil {
		return nil, resp.Header.Get("Content-Type"), resp.StatusCode, err
	}
	if int64(len(body)) > maxBodySize {
		return nil, resp.Header.Get("Content-Type"), resp.StatusCode, fmt.Errorf("image response exceeds %d bytes", maxBodySize)
	}
	if resp.StatusCode != http.StatusOK {
		return body, resp.Header.Get("Content-Type"), resp.StatusCode, nil
	}

	contentType := imageContentType(body, resp.Header.Get("Content-Type"))
	if contentType == "" {
		return nil, resp.Header.Get("Content-Type"), resp.StatusCode, fmt.Errorf("upstream response is not an image")
	}
	return body, contentType, resp.StatusCode, nil
}

func imageContentType(body []byte, upstream string) string {
	contentType := strings.TrimSpace(strings.SplitN(upstream, ";", 2)[0])
	if strings.HasPrefix(strings.ToLower(contentType), "image/") {
		return contentType
	}
	detected := http.DetectContentType(bytes.TrimSpace(body))
	if strings.HasPrefix(strings.ToLower(detected), "image/") {
		return detected
	}
	return ""
}
