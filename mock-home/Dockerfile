# syntax=docker/dockerfile:1

FROM golang:1.19.3-alpine AS base

RUN apk --no-cache add git
RUN apk --no-cache add bash

WORKDIR /app/mock-home

COPY ./mock-home/go.* ./

RUN go mod download

WORKDIR /app

COPY . .

WORKDIR /app/mock-home

RUN chmod +x ./build.sh

RUN ./build.sh

FROM alpine:latest

WORKDIR /app

COPY --from=base /app/mock-home/_output/mock-home ./

EXPOSE 8000

CMD [ "./mock-home" ]

