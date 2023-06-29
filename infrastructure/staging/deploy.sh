#!/bin/bash

export $(grep -v '^#' ../../.staging.env | xargs)

echo $DOCKER_REGISTRY_PASSWORD | docker login --username $DOCKER_REGISTRY_USERNAME --password-stdin

cd ../../.
docker build -t kallepan/pcr-backend:dev -f Dockerfile.staging .
docker push kallepan/pcr-backend:dev

cd infrastructure/staging
kubectl --kubeconfig=$KUBECONFIG delete deployment pcr-backend -n pcr
kubectl --kubeconfig=$KUBECONFIG delete deployment postgres -n pcr 
kubectl --kubeconfig=$KUBECONFIG kustomize . > run.yaml
kubectl --kubeconfig=$KUBECONFIG apply -f run.yaml
