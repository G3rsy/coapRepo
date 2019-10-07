FROM golang:latest

LABEL version="1.0"

RUN mkdir /go/src/server

RUN go get -u github.com/golang/dep/cmd/dep

ADD server.go /go/src/server

WORKDIR /go/src/server

RUN dep init
RUN dep ensure

# Build the Go app
RUN go build server.go

# Expose port 5688, 5689 to the outside world
EXPOSE 5688
EXPOSE 5689

# Command to run the executable
CMD ./server
