#!/bin/bash
set -euo pipefail
cd $(dirname $0)
gcloud workflows deploy million-celebration --source=workflow.yaml --service-account=daily-batch@${PROJECT}.iam.gserviceaccount.com --location=${REGION}

cd ./functions
go test
gcloud functions deploy million-celebration-load \
  --entry-point=Load \
  --runtime go116 --trigger-http --memory 2048MB --timeout 500s --region ${REGION} \
  --set-secrets YOUTUBE_KEY=YOUTUBE_KEY:latest,MILLION_CELEBRATION_SECRETS_JSON=MILLION_CELEBRATION_SECRETS_JSON:latest \
  --set-env-vars PROJECT=${PROJECT},ENV=production

