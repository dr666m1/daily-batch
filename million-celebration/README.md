# 100万再生お祝いbot

## 概要
YouTubeのAPI経由で動画の再生回数を取得し、100万再生達成を検知したら自動でお祝いするbot。
あくまで自分のためのbotなので、誰かに使ってもらうための`README.md`にはなっていない。


## 設定
`./functions/channel_list.py`がYouTubeの検索条件を決めるファイル。書き換えたら`./functions/deploy.sh`を実行して反映すること。

```
# channels
channels = [
    {
        "name": "xxxxx",
        "playlists": ["PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"],
        "tag": "#xxxxx"
    }
]
```

### name
Twitterで投稿する際の表記。敬称も忘れずに。

### playlists
再生リストのURLに`list=xxx`の形式で含まれるもの。複数指定する場合、重複があっても問題ない。

### tag
Twitterで投稿する際のハッシュタグ。複数付ける場合は半角スペースで区切る。

## テーブル定義
```sql
CREATE SCHEMA IF NOT EXISTS million_celebration;
CREATE TABLE IF NOT EXISTS million_celebration.view_count (
  dt DATE,
  playlist_id string,
  video_id string,
  view_count int64
)
PARTITION BY dt
OPTIONS (
  partition_expiration_days = 30
);
CREATE TABLE IF NOT EXISTS million_celebration.view_count_dev (
  dt DATE,
  playlist_id string,
  video_id string,
  view_count int64
)
PARTITION BY dt
OPTIONS (
  partition_expiration_days = 3
);
```
