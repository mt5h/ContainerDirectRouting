# syntax=docker/dockerfile:1

FROM golang:1.19.3-alpine AS base

RUN apk --no-cache add git
RUN apk --no-cache add bash

WORKDIR /app/spawner

COPY ./spawner/go.* ./

RUN go mod download

WORKDIR /app

COPY . .

RUN chmod +x ./spawner/build.sh

RUN ./spawner/build.sh

FROM alpine:latest

WORKDIR /app

COPY --from=base /app/_output/spawner ./

EXPOSE 8008

CMD [ "./spawner" ]

