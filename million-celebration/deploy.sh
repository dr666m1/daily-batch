#!/bin/bash
cd $(dirname $0)
gcloud workflows deploy million-celebration --source=workflow.yaml --service-account=daily-batch@${project}.iam.gserviceaccount.com

cd ./cloud-functions
gcloud functions deploy million-celebration-upload --entry-point main_upload --runtime python37 --trigger-http --memory 2048MB --timeout 500s
gcloud functions deploy million-celebration-tweet --entry-point main_tweet --runtime python37 --trigger-http --memory 2048MB --timeout 500s
