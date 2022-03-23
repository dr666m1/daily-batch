package celebrate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/dr666m1/daily-batch/million-celebration/date"
	"google.golang.org/api/iterator"
)

type ChannelVideoId struct {
	ChannelId string `bigquery:"channel_id"`
	VideoId   string `bigquery:"video_id"`
}

type Secrets struct {
	ConsumerKey    string `json:"CONSUMER_KEY"`
	ConsumerSecret string `json:"CONSUMER_SECRET"`
	Token          string `json:"TOKEN"`
	TokenSecret    string `json:"TOKEN_SECRET"`
}

func createMessage(name string, threshold int, video_id string, tags []string) string {
	var hashtags []string
	for _, tag := range tags {
		hashtags = append(hashtags, "#"+tag)
	}
	msg := fmt.Sprintf(
		"%v%d万再生おめでとう！\nhttps://www.youtube.com/watch?v=%v\n%v",
		name,
		threshold,
		video_id,
		strings.Join(hashtags, " "),
	)
	return msg
}

func newClient() (*twitter.Client, error) {
	var secrets Secrets
	err := json.Unmarshal(
		[]byte(os.Getenv("MILLION_CELEBRATION_SECRETS_JSON")),
		&secrets,
	)
	if err != nil {
		return nil, err
	}
	config := oauth1.NewConfig(secrets.ConsumerKey, secrets.ConsumerSecret)
	token := oauth1.NewToken(secrets.Token, secrets.TokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)
	return client, nil
}

func queryChannelVideoIds(client *bigquery.Client, today string, threshold int) ([]ChannelVideoId, error) {
	query := `
WITH
  yesterday AS (
    SELECT channel_id, video_id, MIN(view_count) AS view_count_yesterday
    FROM million_celebration.%v
    WHERE dt = DATE '%v'
    GROUP BY 1, 2
  ),
  today AS (
		SELECT channel_id, video_id, MIN(view_count) AS view_count_today
    FROM million_celebration.%v
    WHERE dt = DATE '%v'
    GROUP BY 1, 2
  )
SELECT DISTINCT channel_id, video_id
FROM yesterday INNER JOIN today USING(channel_id, video_id)
WHERE
  view_count_yesterday < %v
  AND %v <= view_count_today
ORDER BY 1, 2
`
	table := "view_count_dev"
	if os.Getenv("ENV") == "production" {
		table = "view_count"
	}

	q := client.Query(fmt.Sprintf(query, table, date.OneDayBefore(today), table, today, threshold, threshold))
	ctx := context.Background()
	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}

	var rows []ChannelVideoId
	for {
		var row ChannelVideoId
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func Tweet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	bqClient, err := bigquery.NewClient(ctx, os.Getenv("PROJECT"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	twitterClient, err := newClient()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	thresholds := []int{100} // unit: 10_000
	for _, threshold := range thresholds {
		channelVideoIds, err := queryChannelVideoIds(bqClient, date.Today(), threshold*10_000)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		for _, cv := range channelVideoIds {
			var name string
			var tags []string
			for _, channel := range channels {
				if channel.channelId == cv.ChannelId {
					name = channel.name
					tags = channel.tags
				}
			}
			if name == "" || len(tags) == 0 {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Cannot get the information for the channel: %v", cv.ChannelId)))
				return
			}
			msg := createMessage(name, threshold, cv.VideoId, tags)
			twitterClient.Statuses.Update(msg, nil)
		}
	}
}
