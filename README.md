# Subscriptions POC

### Local Setup


##### First Time?
```
  brew install tilt
  brew install k3d
  brew install kubectx 
  brew install istioctl
  
  mkdir -p /tmp/k3dvol/
  k3d cluster create factoryCluster --volume /tmp/k3dvol:/tmp/k3dvol --registry-create local-factory-registry
  kubectl create ns subscriptions
  kubectl label namespace subscriptions istio-injection=enabled
```

##### To run locally in K3d:
```
  ./run.sh
```

##### To run locally in docker
```
  ./run-docker.sh
```
