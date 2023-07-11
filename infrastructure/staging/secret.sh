#!/bin/bash
#
# Create Kubernetes Secrets from local secrets.env
#

export $(grep -v '^#' ../../.staging.env | xargs)

echo ${DOCKER_REGISTRY_SERVER}

kubectl --kubeconfig=$KUBECONFIG delete secret oe-secrets -n pcr
kubectl --kubeconfig=$KUBECONFIG create secret generic oe-secrets --from-env-file=../../.staging.env -n pcr
kubectl --kubeconfig=$KUBECONFIG create secret docker-registry regcred --docker-server=${DOCKER_REGISTRY_SERVER} --docker-username=${DOCKER_REGISTRY_USERNAME} --docker-password=${DOCKER_REGISTRY_PASSWORD} -n pcr