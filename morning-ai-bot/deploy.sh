#!/bin/bash
cd $(dirname $0)
gcloud workflows deploy morning-ai-bot --source=workflow.yaml --service-account=daily-batch@${project}.iam.gserviceaccount.com

cd ./cloud-functions
gcloud functions deploy morning-ai-bot --entry-point morning_ai_bot --runtime python37 --trigger-http --memory 2048MB --timeout 500s
