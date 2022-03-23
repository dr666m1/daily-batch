package celebrate

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/google/go-cmp/cmp"
)

func Test_createMessage(t *testing.T) {
	expected := `xxxxxちゃん100万再生おめでとう！
https://www.youtube.com/watch?v=yyyyy
#aaaaa #bbbbb`
	actual := createMessage("xxxxxちゃん", 100, "yyyyy", []string{"aaaaa", "bbbbb"})
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

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
		expected  []ChannelVideoId
	}{
		{"100", 100, []ChannelVideoId{
			{"xxxxx", "a"},
		}},
		{"200", 200, nil},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			rows, err := queryChannelVideoIds(client, "2000-01-02", d.threshold)
			if err != nil {
				t.Error(err.Error())
			}
			if diff := cmp.Diff(d.expected, rows); diff != "" {
				t.Error(diff)
			}
		})
	}
}
