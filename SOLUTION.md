# Solution

## Setup
1. `git clone git@github.com:tutti-ch/backend-coding-task-2024-11-milan-keca.git`
2. `go mod download`
3. `mkdir /tmp/smg/` - create a tmp directory for storing images
4. `export BASE_PATH="/tmp/smg/"`

## Starting the server
1. `go run main.go`

## Running tests
1. `mkdir /tmp/smgtest`
2. `go test ./...`

## Usage
1. `curl -v -F image=@IMAGE_PATH SERVICE_URL`
  * `curl -v -F image=@testdata/testimage_small.jpg localhost:3000/upload`

## Running with Kubernetes
Requirements:
1. minikube
2. kubectl
3. docker

### Resources
1. Deployment for the api
2. Service
3. ConfigMap

### Setup
1. `docker build -t tuti .`
2. `minikube image load tuti:latest`
    * Load image to minikube's registry
    * Image pull policy is set to Never
3. `minikube start`
4. `kubectl apply -f manifests/`

To validate everything is set up run:
`kubectl get pod,svc,configMap`

The api is exposed via a NodePort on port `30005`

### Interacting with the cluster
1. `minikube service api-service --url` - Find the address of your service
2. `curl -v -F image=@IMAGE_PATH SERVICE_URL`

To check the logs of the pods run:
`kubectl logs POD_NAME`

To check if the file is stored on the container's storage run:
1. `k exec --stdin --tty POD_NAME -- sh`
2. `du -h FILEPATH` - check size of the image
