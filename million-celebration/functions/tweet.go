package celebrate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/dr666m1/daily-batch/million-celebration/date"
	"google.golang.org/api/iterator"
)

type PlaylistVideoId struct {
	PlaylistId string `bigquery:"playlist_id"`
	VideoId    string `bigquery:"video_id"`
}

type Secrets struct {
	ConsumerKey    string `json:"CONSUMER_KEY"`
	ConsumerSecret string `json:"CONSUMER_SECRET"`
	Token          string `json:"TOKEN"`
	TokenSecret    string `json:"TOKEN_SECRET"`
}

func Tweet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	_, err := bigquery.NewClient(ctx, os.Getenv("PROJECT"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
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

func queryPlaylistVideoIds(client *bigquery.Client, today string, threshold int) ([]PlaylistVideoId, error) {
	query := `
WITH
  yesterday AS (
    SELECT playlist_id, video_id, MIN(view_count) AS view_count_yesterday
    FROM million_celebration.%v
    WHERE dt = DATE '%v'
    GROUP BY 1, 2
  ),
  today AS (
		SELECT playlist_id, video_id, MIN(view_count) AS view_count_today
    FROM million_celebration.%v
    WHERE dt = DATE '%v'
    GROUP BY 1, 2
  )
SELECT DISTINCT playlist_id, video_id
FROM yesterday INNER JOIN today USING(playlist_id, video_id)
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

	var rows []PlaylistVideoId
	for {
		var row PlaylistVideoId
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
