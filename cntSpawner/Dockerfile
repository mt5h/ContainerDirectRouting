# syntax=docker/dockerfile:1

FROM golang:1.19.3-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY *.go ./

RUN go build -o /cntSpawner

EXPOSE 8008

CMD [ "/cntSpawner" ]

