#!/bin/bash
set -euo pipefail
cd $(dirname $0)
export PROJECT=$(gcloud config get-value project)

for x in *; do
  if [ -d "$x" ] && [ -x "${x}/deploy.sh" ]; then
    ./${x}/deploy.sh
  fi
done
gcloud workflows deploy daily-batch --source=workflow.yaml --service-account=daily-batch@${project}.iam.gserviceaccount.com --location=${REGION}
