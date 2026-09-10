package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"openfeed/internal/model"
	"openfeed/internal/parser"
	"openfeed/internal/provider"
	"openfeed/internal/telemirror"
)

const defaultFeedLimit = 50
const maxFeedLimit = 100

func Channel(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/channel/")
	limit := defaultFeedLimit
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > maxFeedLimit {
		limit = maxFeedLimit
	}

	before := 0
	if raw := r.URL.Query().Get("before"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			before = n
		}
	}

	channel, err := fetchChannelPage(name, before, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(channel)
}

func fetchChannelPage(name string, beforeID, limit int) (*model.Channel, error) {
	var (
		result  *model.Channel
		seen    = map[string]struct{}{}
		cursor  = beforeID
		hasMore = false
	)

	for lenPosts := 0; lenPosts < limit; {
		var html string
		var err error
		if cursor > 0 {
			html, err = provider.Default.ClientFetchHTMLBefore(name, cursor)
		} else {
			html, err = provider.Default.LoadChannel(name)
		}
		if err != nil {
			return nil, err
		}

		ch, posts, err := telemirror.ParseHTML(string(html))
		if err != nil {
			return nil, err
		}
		converted := parser.Convert(ch, posts)
		if result == nil {
			result = converted
		} else {
			for _, post := range converted.Posts {
				if _, ok := seen[post.ID]; ok {
					continue
				}
				result.Posts = append(result.Posts, post)
			}
		}
		for _, post := range converted.Posts {
			seen[post.ID] = struct{}{}
		}

		lenPosts := len(result.Posts)
		if len(converted.Posts) == 0 {
			break
		}
		if lenPosts >= limit {
			hasMore = true
			break
		}
		if len(converted.Posts) < 10 {
			break
		}

		oldest := result.Posts[len(result.Posts)-1].ID
		parts := strings.Split(oldest, "/")
		if len(parts) == 0 {
			break
		}
		next, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil || next <= 0 || next == cursor {
			break
		}
		cursor = next
	}

	if result == nil {
		result = &model.Channel{Posts: []model.Post{}}
	}
	if len(result.Posts) > limit {
		result.Posts = result.Posts[:limit]
	}
	result.HasMore = hasMore
	if len(result.Posts) > 0 {
		parts := strings.Split(result.Posts[len(result.Posts)-1].ID, "/")
		if len(parts) > 0 {
			if next, err := strconv.Atoi(parts[len(parts)-1]); err == nil && next > 0 {
				result.NextBefore = strconv.Itoa(next)
			}
		}
	}
	return result, nil
}
