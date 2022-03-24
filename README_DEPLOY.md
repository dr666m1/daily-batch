# Deploy
## サービスアカウント
daily-batchという名称のサービスアカウントを作成する。
必要な権限は以下。

* Cloud Functions Invoker
* Logs Writer
* Secret Manager Secret Accessor
* Workflows Invoker

## Cloud Scheduler
Workflowsのdaily-batchを実行する設定を行う（[参考](https://cloud.google.com/workflows/docs/schedule-workflow)）。

## Secret Manager
Cloud Functionsに権限を付与しておく（[参考](https://cloud.google.com/functions/docs/configuring/secrets)）。
そして以下のSECRETを作成する。

* YOUTUBE\_KEY
* TWITTER\_TOKEN
* MILLION\_CELEBRATION\_SECRETS\_JSON（必要なキーは以下）
  * CONSUMER\_KEY
  * CONSUMER\_SECRET
  * TOKEN
  * TOKEN\_SECRET
* LINE\_TOKEN\_SANDBOX

## .envrc
以下のように環境変数を設定し`deploy.sh`を実行。

```bash
# .envrc
## デプロイ先の指定
export PROJECT=xxxxx
export REGION=yyyyy

## テスト用
export YOUTUBE_KEY=zzzzz
export MILLION_CELEBRATION_SECRETS_JSON='{"CONSUMER_KEY": "aaaaa", "CONSUMER_SECRET": "bbbbb", "TOKEN": "ccccc", "TOKEN_SECRET": "ddddd"}'
```
