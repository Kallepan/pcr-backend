#!/bin/bash

# copy staging migration to migrations
export MIGRATION_FILE="00010_staging.up.sql"
cp $MIGRATION_FILE migrations/$MIGRATION_FILE

export $(grep -v '^#' .staging.env | xargs)

echo $DOCKER_REGISTRY_PASSWORD | docker login --username $DOCKER_REGISTRY_USERNAME --password-stdin

kubectl create namespace $NAMESPACE --kubeconfig=$KUBECONFIG

kubectl delete secret secrets -n $NAMESPACE --kubeconfig=$KUBECONFIG
kubectl create secret generic secrets \
    --from-env-file=.staging.env \
    -n $NAMESPACE \
    --kubeconfig=$KUBECONFIG

kubectl delete secret regcred -n $NAMESPACE --kubeconfig=$KUBECONFIG
kubectl create secret docker-registry -n $NAMESPACE \
    regcred \
    --docker-server=$DOCKER_REGISTRY_SERVER \
    --docker-username=$DOCKER_REGISTRY_USERNAME \
    --docker-password=$DOCKER_REGISTRY_PASSWORD \
    --kubeconfig=$KUBECONFIG

# Delete old deployments
kubectl delete deployment $DOCKER_REGISTRY_REPOSITORY -n $NAMESPACE --kubeconfig=$KUBECONFIG
kubectl delete deployment postgres -n $NAMESPACE --kubeconfig=$KUBECONFIG

# build and push staging image
docker build -t $DOCKER_REGISTRY_USERNAME/$DOCKER_REGISTRY_REPOSITORY:${VERSION} .
docker push $DOCKER_REGISTRY_USERNAME/$DOCKER_REGISTRY_REPOSITORY:${VERSION}

# Deploy new staging image
cd infrastructure/staging
kubectl kustomize . > run.yaml
sed -i "s/IMAGE_TAG/${VERSION}/g" run.yaml
kubectl apply -f run.yaml --kubeconfig=$KUBECONFIG

# remove staging migration from migrations
cd ../../
rm migrations/$MIGRATION_FILE