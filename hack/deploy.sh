#!/bin/bash
#
# Deploy to production
#
export $(grep -v '^#' .prod.env | xargs)

echo $DOCKER_REGISTRY_PASSWORD | docker login --username $DOCKER_REGISTRY_USERNAME --password-stdin

kubectl delete secret secrets -n $NAMESPACE --kubeconfig=$KUBECONFIG
kubectl create secret generic secrets \
    --from-env-file=.prod.env \
    -n $NAMESPACE \
    --kubeconfig=$KUBECONFIG

kubectl delete secret regcred -n $NAMESPACE --kubeconfig=$KUBECONFIG
kubectl create secret docker-registry -n $NAMESPACE \
    regcred \
    --docker-server=$DOCKER_REGISTRY_SERVER \
    --docker-username=$DOCKER_REGISTRY_USERNAME \
    --docker-password=$DOCKER_REGISTRY_PASSWORD \
    --kubeconfig=$KUBECONFIG

# build and push production image
docker build --platform linux/amd64 -t $DOCKER_REGISTRY_USERNAME/$DOCKER_REGISTRY_REPOSITORY:${VERSION} -f Dockerfile.prod .
docker push $DOCKER_REGISTRY_USERNAME/$DOCKER_REGISTRY_REPOSITORY:${VERSION}

# Deploy new production image
cd infrastructure/prod
kubectl kustomize . > run.yaml
sed -i '' "s/IMAGE_TAG/${VERSION}/g" run.yaml
kubectl apply -f run.yaml --kubeconfig=$KUBECONFIG