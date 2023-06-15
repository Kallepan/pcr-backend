#!/bin/bash
#
# Create Kubernetes Secrets from local secrets.env
#

export $(grep -v '^#' ../../.staging.env | xargs)

echo ${DOCKER_REGISTRY_SERVER}

kubectl create secret generic secrets --from-env-file=../../.staging.env -n pcr
kubectl create secret docker-registry regcred --docker-server=${DOCKER_REGISTRY_SERVER} --docker-username=${DOCKER_REGISTRY_USERNAME} --docker-password=${DOCKER_REGISTRY_PASSWORD} -n pcr