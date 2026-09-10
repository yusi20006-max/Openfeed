package parser

import (
	"net/url"
	"sort"
	"time"

	"openfeed/internal/model"
	"openfeed/internal/telemirror"
)

// proxyImage rewrites a translate.goog image URL to go through our own
// /api/image endpoint, so the browser fetches it from our server (which
// has unrestricted network access) instead of connecting to
// *.translate.goog directly — a host blocked on many Iranian networks
// without a VPN, unlike the channel HTML which is fetched server-side.
func proxyImage(raw string) string {
	if raw == "" {
		return raw
	}
	return "/api/image?u=" + url.QueryEscape(raw)
}

// Convert maps the rich telemirror result (real parser: media, replies,
// forwards, sanitized HTML) onto openfeed's own simpler model used by the frontend.
func Convert(ch *telemirror.Channel, posts []telemirror.Post) *model.Channel {

	dst := &model.Channel{
		Title:       ch.Title,
		Username:    ch.Username,
		Avatar:      proxyImage(ch.Photo),
		Description: ch.Description,
		Subscribers: ch.Subscribers,
		Posts:       make([]model.Post, 0, len(posts)),
	}

	// Keep pagination/load-more semantics intact: only establish the order of
	// the posts already returned by the provider. Newest posts are first, while
	// older posts remain in their relative order and can still be appended by a
	// future page without reversing the accumulated feed.
	sort.SliceStable(posts, func(i, j int) bool {
		if posts[i].Time.IsZero() {
			return false
		}
		if posts[j].Time.IsZero() {
			return true
		}
		return posts[i].Time.After(posts[j].Time)
	})

	for _, p := range posts {

		post := model.Post{
			ID:     p.ID,
			Author: p.Author,
			Text:   p.Text,
			Views:  p.Views,
		}

		if !p.Time.IsZero() {
			post.Date = p.Time.Format(time.RFC3339)
		}

		for _, m := range p.Media {

			entry := model.Media{
				Type:     m.Type,
				Ratio:    m.Ratio,
				Download: m.Download,
				Duration: m.Duration,
				Title:    m.Title,
				Subtitle: m.Subtitle,
			}

			// Thumb is the actual image URL, rewritten by proxyImage to
			// go through our own /api/image endpoint so it loads even
			// when Telegram's CDN and translate.goog are both blocked
			// on the client's network.
			if m.Thumb != "" {
				entry.URL = proxyImage(m.Thumb)
			} else if m.Download == "" {
				// Nothing displayable and nothing downloadable (e.g. a
				// poll, or a media kind we don't parse) — skip it.
				continue
			}

			post.Media = append(post.Media, entry)
		}

		dst.Posts = append(dst.Posts, post)

	}

	return dst

}
