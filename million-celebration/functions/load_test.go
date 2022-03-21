package celebrate

import (
	"context"
	"os"
	"testing"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func Test_callPlaylistItems(t *testing.T) {
	ctx := context.Background()
	client, err := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("YOUTUBE_KEY")))
	if err != nil {
		t.Error(err.Error())
	}
	videos, token, err := callPlaylistItems(
		client,
		"PL0bHKk6wuUGIAmzzqdVMynRrAOi8odYFQ",
		"",
	)
	if err != nil {
		t.Error(err.Error())
	}
	if len(videos) == 0 {
		t.Error("Got empty response")
	}
	if token == "" {
		t.Error("Next page token cannot be empty because this playlist is loger than 50.")
	}
}

func Test_callVideos(t *testing.T) {
	ctx := context.Background()
	client, err := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("YOUTUBE_KEY")))
	if err != nil {
		t.Error(err.Error())
	}

	data := []struct {
		name      string
		videoIds  []string
		batchSize int
	}{
		{"one id", []string{"On5gvNPjXGQ"}, 50},
		{"two ids", []string{"On5gvNPjXGQ", "6iYKlBSQu1g"}, 50},
		{"tree ids", []string{"On5gvNPjXGQ", "6iYKlBSQu1g", "RvB-kv9q7Pk"}, 2},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			videos, err := callVideos(client, d.videoIds, d.batchSize)
			if err != nil {
				t.Fatal(d.name, err.Error())
			}
			if len(videos) != len(d.videoIds) {
				t.Errorf("incorrect length: expected %d, got %d", len(d.videoIds), len(videos))
			}
			if videos[0].viewCount < 8_000_000 {
				t.Errorf("incorrect view count: expected more than 8M, got %d", videos[0].viewCount)
			}
		})
	}
}

func Test_insertRows(t *testing.T) {
	data := []struct {
		name  string
		rows  []Row
		valid bool
	}{
		{"valid", []Row{{"2020-01-01", "playlist_id", "video_id", 1}}, true},
		{"invalid", []Row{{"abc", "playlist_id", "video_id", 1}}, false},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err := insertRows(d.rows)
			if d.valid && err != nil {
				t.Error(d.name, err.Error())
			} else if !d.valid && err == nil {
				t.Errorf("%v unexpectedly successed.", d.name)
			}
		})
	}
}
