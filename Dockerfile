FROM golang:latest

WORKDIR /go/src/project1/

ADD go.mod go.work go.work.sum main.go ./
ADD service handlers repository ./

EXPOSE 8080

CMD ["go", "run", "main.go"]