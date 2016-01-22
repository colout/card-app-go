#!/bin/bash

# Stop golang container (will error if not started
docker stop go

# Delete it...
docker rm go

# Run start golang container
docker run -d --name go -p 127.0.0.1:8080:8080 -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.3 /bin/bash -c "go get -t; go run -v main.go"
