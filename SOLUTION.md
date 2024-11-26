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
