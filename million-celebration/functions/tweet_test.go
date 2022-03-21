package celebrate

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/google/go-cmp/cmp"
)

func Test_newClient(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Error(err.Error())
	}
	client.Statuses.Update("test", nil)
}

func Test_queryApplicableVideos(t *testing.T) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, os.Getenv("PROJECT"))
	if err != nil {
		t.Error(err.Error())
	}

	data := []struct {
		name      string
		threshold int
		expected  []PlaylistVideoId
	}{
		{"100", 100, []PlaylistVideoId{
			{"PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "a"},
		}},
		{"200", 200, nil},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			rows, err := queryPlaylistVideoIds(client, "2000-01-02", d.threshold)
			if err != nil {
				t.Error(err.Error())
			}
			if diff := cmp.Diff(d.expected, rows); diff != "" {
				t.Error(diff)
			}
		})
	}
}
