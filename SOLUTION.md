# Solution
## Implementation details
### Request Lifetime
When a request is received, bytes are read from the payload and the required
checks are done.

If the checks are ok:
1. Image id is generated
2. A message is sent to the workers
3. 200 is returned

If the checks are not ok:
1. Proper errors are returned

### Communication between the api and workers
Workers are started in separate go routines and receive messages through a
blocking channel (unbuffered). The reason for this is to be able to handle
timeouts easily.
Timeouts are handled by using a select statement with:
1. sending a message to the channel
2. timeing out after 100ms
Since the channel is blocking, the message won't be sent until a worker is ready
to receive a message.

### Passing the payload to the workers
The images are passed as bytes to the workers.
This is not ideal in terms of memory, but since the number of workers are
limited, the decision was made to work with memory instead of storing the images
in temp storage and cleaning it up later.

### Graceful shutdown
A struct is made to handle both the server and the workers, so that the
entrypoint stays minimal. This allows easier testing of most of the
functionality.
The api and workers are all started in separate go routines.
The main go routine is waiting for a signal from the OS (SIGINT or SIGTERM).
Once the signal is passed to the main go routine, it will gracefully shut down
the server, and also close the channel for the workers. Workers are waited to
finish their jobs, and the main go routine exits.

### Additional improvements
1. Server should parametrized with port number
2. Request context could be passed to the workers for easier monitoring
3. Testing can be refactored and cleaned
    * Some errors are not handled like reading an image to test
    * All of the application specific errors are handled
    * Most of the application is tested
4. Naming can be better in some places

## Setup
1. `git clone git@github.com:tutti-ch/backend-coding-task-2024-11-milan-keca.git`
2. `go mod download`
3. `mkdir /tmp/smg/` - create a tmp directory for storing images
4. `export BASE_PATH="/tmp/smg/"`
5. `export N_WORKERS=N` - if not set, number of cpus will be used

## Starting the server
1. `go run cmd/api/server.go`

Additionally, to simulate many requests run:
1. `go run cmd/simulation/simulate.go`
This will send thousands of requests to the api simultaneusly, to catch timeout
exceptions use a smaller number of workers (3-5).

## Running tests
1. `go test ./...`

## Usage
1. `curl -v -F image=@IMAGE_PATH SERVICE_URL`
  * `curl -v -F image=@testdata/testimage_small.jpg localhost:3000/upload`

## Making it production ready
1. Optimize docker image
2. Create a message queue and decouple workers from the api
3. Communicate through the message queue
4. Store each job request and job run
5. Using an object store to upload result images
6. Proper logging and log aggregation
    * propagate ID of each request to the workers for easier debugging and monitoring
7. API rate limiting
8. Provide endpoints for viewing details about running or finished jobs
9. CICD pipelines
    * running tests
    * uploading images to registry
    * triggering deployments
10. ...

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
    * to validate the image is loaded, run `minikube image ls`
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
1. `kubectl exec --stdin --tty POD_NAME -- sh`
2. `du -h FILEPATH` - check size of the image
