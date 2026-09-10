package parser

import (
	"testing"
	"time"

	"openfeed/internal/telemirror"
)

func TestConvertOrdersPostsNewestFirst(t *testing.T) {
	newest := time.Date(2026, 9, 10, 12, 0, 0, 0, time.UTC)
	middle := newest.Add(-time.Hour)
	oldest := newest.Add(-2 * time.Hour)

	posts := []telemirror.Post{
		{ID: "test/1", Time: oldest},
		{ID: "test/3", Time: newest},
		{ID: "test/2", Time: middle},
	}

	channel := Convert(&telemirror.Channel{Username: "test"}, posts)
	if got, want := len(channel.Posts), 3; got != want {
		t.Fatalf("post count = %d, want %d", got, want)
	}

	wantIDs := []string{"test/3", "test/2", "test/1"}
	for i, want := range wantIDs {
		if channel.Posts[i].ID != want {
			t.Fatalf("post[%d].ID = %q, want %q", i, channel.Posts[i].ID, want)
		}
	}
}

func TestConvertKeepsImageProxyLocal(t *testing.T) {
	posts := []telemirror.Post{{
		ID: "test/1",
		Time: time.Now(),
		Media: []telemirror.Media{{Type: "photo", Thumb: "https://example.translate.goog/image.jpg?_x_tr_sl=auto"}},
	}}

	channel := Convert(&telemirror.Channel{Username: "test"}, posts)
	if got := channel.Posts[0].Media[0].URL; got == "" || got[:len("/api/image?u=")] != "/api/image?u=" {
		t.Fatalf("image URL = %q, want local /api/image proxy", got)
	}
}
