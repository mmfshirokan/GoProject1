FROM golang:latest

WORKDIR /projects/project1/GoProject1

ADD go.mod go.sum main.go ./
ADD internal /projects/project1/GoProject1/internal

EXPOSE 8080

CMD ["go", "run", "main.go"]