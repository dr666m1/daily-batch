package celebrate

import (
	"context"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/dr666m1/daily-batch/million-celebration/date"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Row struct {
	Dt         string `bigquery:"dt"`
	channelId  string `bigquery:"channel_id"`
	PlaylistId string `bigquery:"playlist_id"`
	VideoId    string `bigquery:"video_id"`
	ViewCount  int    `bigquery:"view_count"`
}

type videoResponse struct {
	videoId   string
	viewCount int
}

func callPlaylistItems(client *youtube.Service, playlistId string, nextPageToken string) ([]string, string, error) {
	call := client.PlaylistItems.List([]string{"snippet", "status"})
	call.MaxResults(50)
	call.PlaylistId(playlistId)
	if nextPageToken != "" {
		call.PageToken(nextPageToken)
	}
	resp, err := call.Do()
	if err != nil {
		return []string{}, "", err
	}

	var videoIds []string
	for _, item := range resp.Items {
		if item.Status.PrivacyStatus != "public" { // unlisted, private
			continue
		}
		videoIds = append(videoIds, item.Snippet.ResourceId.VideoId)
	}

	return videoIds, resp.NextPageToken, err
}

func callVideos(client *youtube.Service, videoIds []string, batchSize int) ([]videoResponse, error) {
	var start int
	var items []*youtube.Video

	for start < len(videoIds) {
		end := min(start+batchSize, len(videoIds))
		ids := videoIds[start:end]
		start = end
		call := client.Videos.List([]string{"statistics"})
		call.Id(strings.Join(ids, ","))
		resp, err := call.Do()
		if err != nil {
			return nil, err
		}
		items = append(items, resp.Items...)
	}

	var videos []videoResponse
	for _, item := range items {
		videos = append(videos, videoResponse{item.Id, int(item.Statistics.ViewCount)})
	}
	return videos, nil
}

func insertRows(rows []Row) error {
	ctx := context.Background()
	bqClient, err := bigquery.NewClient(ctx, os.Getenv("PROJECT"))
	if err != nil {
		return err
	}
	defer bqClient.Close()

	dataset := bqClient.Dataset("million_celebration")
	table := dataset.Table("view_count_dev")
	if env := os.Getenv("ENV"); env == "production" {
		table = dataset.Table("view_count")
	}
	inserter := table.Inserter()
	if err := inserter.Put(ctx, rows); err != nil {
		return err
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Load(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	youtubeClient, err := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("YOUTUBE_KEY")))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	dt := date.Today()
	var rows []Row
	for _, c := range channels {
		for _, pl := range c.playlists {
			var nextPageToken string
			var videoIds []string
			for {
				ids, token, err := callPlaylistItems(youtubeClient, pl, nextPageToken)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
				videoIds = append(videoIds, ids...)
				if token == "" {
					break
				}
				nextPageToken = token
			}
			videos, err := callVideos(youtubeClient, videoIds, 50)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			for _, v := range videos {
				rows = append(rows, Row{dt, c.channelId, pl, v.videoId, v.viewCount})
			}
		}
	}

	if err := insertRows(rows); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("done!"))
}
