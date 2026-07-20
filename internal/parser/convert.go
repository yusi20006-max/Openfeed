package parser

import (
	"time"

	"openfeed/internal/model"
	"openfeed/internal/telemirror"
)

// Convert maps the rich telemirror result (real parser: media, replies,
// forwards, sanitized HTML) onto openfeed's own simpler model used by
// the frontend.
func Convert(ch *telemirror.Channel, posts []telemirror.Post) *model.Channel {

	dst := &model.Channel{
		Title:       ch.Title,
		Username:    ch.Username,
		Avatar:      ch.Photo,
		Description: ch.Description,
		Subscribers: ch.Subscribers,
		Posts:       make([]model.Post, 0, len(posts)),
	}

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

			// Thumb is the actual image URL served back through the
			// translate.goog proxy, so it loads even when Telegram's own
			// CDN is blocked.
			if m.Thumb != "" {
				entry.URL = m.Thumb
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
