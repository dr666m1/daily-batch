# 100万再生お祝いbot

## 概要
YouTubeのAPI経由で動画の再生回数を取得し、100万再生達成を検知したら自動でお祝いするbot。
あくまで自分のためのbotなので、誰かに使ってもらうための`README.md`にはなっていない。

## 設定
`./functions/channels.go`がYouTubeの検索条件を決めるファイル（Gitの管理外とした）。

```go
package celebrate

type channel struct {
	channelId string
	name      string
	playlists []string
	tags      []string
}

var channels = []channel{
	{
		id: "xxxxx",
		name: "xxxxx",
		playlists: []string{
			"PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		},
		tags: []string{"xxxxx"},
	},
}
```

### id
`https://www.youtube.com/channel/xxxxx`の末尾。
一意になる識別子がほしいだけだから本当は何でもいい。

### name
Twitterで投稿する際の表記。敬称も忘れずに。

### playlists
再生リストのURLに`list=xxx`の形式で含まれるもの。複数指定する場合、重複があっても問題ない。

### tags
Twitterで投稿する際のハッシュタグ。`#`は不要。

## テーブル定義
```sql
CREATE SCHEMA IF NOT EXISTS million_celebration;
CREATE TABLE IF NOT EXISTS million_celebration.view_count (
  dt DATE,
  channel_id string,
  playlist_id string,
  video_id string,
  view_count int64
)
PARTITION BY dt
OPTIONS (
  partition_expiration_days = 30
);
CREATE OR REPLACE TABLE million_celebration.view_count_dev
LIKE million_celebration.view_count
OPTIONS (
  partition_expiration_days = NULL
);
INSERT million_celebration.view_count_dev VALUES
  (DATE '2000-01-01', 'xxxxx', 'PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx', 'a', 99),
  (DATE '2000-01-01', 'xxxxx', 'PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx', 'a', 100),
  (DATE '2000-01-01', 'xxxxx', 'PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx', 'b', 99),
  (DATE '2000-01-01', 'xxxxx', 'PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx', 'c', 99),
  (DATE '2000-01-02', 'xxxxx', 'PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx', 'a', 100),
  (DATE '2000-01-02', 'xxxxx', 'PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx', 'b', 99),
  (DATE '2000-01-01', 'xxxxx', 'PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx', 'd', 200)
;
```
