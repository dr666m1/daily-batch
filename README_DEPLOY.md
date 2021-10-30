# Dploy
## 準備
必要に応じてCloud SDKのデフォルトregionを設定しておくこと。

```bash
region=xxxxx # us-west1, us-east1, us-central1...
gcloud config set workflows/location $region
gcloud config set functions/region $region
```
### サービスアカウント
daily-batchという名称のサービスアカウントを作成する。必要な権限は以下。

* Logs Writer
* Cloud Functions Invoker
* Workflows Invoker

### Cloud Scheduler
Workflowsのdaily-batchを実行する設定を行う（[参考](https://cloud.google.com/workflows/docs/schedule-workflow)）。

### Secret Manager
Cloud Functionsに権限を付与しておく（[参考](https://cloud.google.com/functions/docs/configuring/secrets)）。
そして以下のSECRETを作成する。

* YOUTUBE\_KEY
* TWITTER\_TOKEN
* MILLION\_CELEBRATION\_SECRETS\_JSON（必要なキーは以下）
  * CONSUMER\_KEY
  * CONSUMER\_SECRET
  * TOKEN
  * TOKEN\_SECRET
* MORNING\_AI\_BOT\_TOKEN
* LINE\_TOKEN\_SANDBOX

## 本番反映
```bash
./deploy.sh
```
