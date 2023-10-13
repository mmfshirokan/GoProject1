FROM golang:latest

ADD go.mod /go/src/echoServer/
ADD go.sum /go/src/echoServer/
ADD echoServer.go /go/src/echoServer/

WORKDIR /go/src/echoServer/

CMD [ "go", "run", "echoServer.go" ]