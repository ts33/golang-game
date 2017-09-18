# Use an official Golang runtime as a base image
FROM golang:1.9

# copy all files in the current directory to the container
ADD . /go/src/github.com/ts33/golang-game

# downloads all dependencies
RUN go get -d ./...

# installs the app
RUN go install github.com/ts33/golang-game

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/golang-game

# Document that the service listens on port 8080.
EXPOSE 3000
