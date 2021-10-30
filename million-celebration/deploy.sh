#!/bin/bash
cd $(dirname $0)
gcloud workflows deploy million-celebration --source=workflow.yaml --service-account=daily-batch@${project}.iam.gserviceaccount.com

cd ./cloud-functions
gcloud beta functions deploy million-celebration-upload \
  --entry-point main_upload --runtime python37 --trigger-http --memory 2048MB --timeout 500s \
  --set-secrets YOUTUBE_KEY=YOUTUBE_KEY:latest,MILLION_CELEBRATION_SECRETS_JSON=MILLION_CELEBRATION_SECRETS_JSON:latest
gcloud beta functions deploy million-celebration-tweet \
  --entry-point main_tweet --runtime python37 --trigger-http --memory 2048MB --timeout 500s \
  --set-secrets YOUTUBE_KEY=YOUTUBE_KEY:latest,MILLION_CELEBRATION_SECRETS_JSON=MILLION_CELEBRATION_SECRETS_JSON:latest

