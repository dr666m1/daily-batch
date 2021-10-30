# Daily Batch
## 概要
GCPの無料枠で構築したオレオレバッチ処理。

## 事前準備
`daily-batch`という名称のサービスアカウントを作成する。必要な権限は以下。

* Logs Writer
* Cloud Functions Invoker
* Workflows Invoker

必要に応じてCloud SDKのデフォルトregionを設定しておく。

```bash
region=xxxxx # us-west1, us-east1, us-central1...
gcloud config set workflows/location $region
gcloud config set functions/region $region
```

Cloud Schedulerから`daily-batch`を実行する設定を行う（[参考](https://cloud.google.com/workflows/docs/schedule-workflow)）。

## 本番反映

```bash
./deploy.sh
```

## TODO
* Cloud BuildまたはGithub Actionsで本番反映の自動化
* LINEにエラー通知
