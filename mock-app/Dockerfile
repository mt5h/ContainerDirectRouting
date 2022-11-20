# syntax=docker/dockerfile:1

FROM golang:1.19.3-alpine AS base

RUN apk --no-cache add git
RUN apk --no-cache add bash

WORKDIR /app/mock-app

COPY ./mock-app/go.* ./

RUN go mod download

WORKDIR /app

COPY . .

WORKDIR /app/mock-app

RUN chmod +x ./build.sh

RUN ./build.sh

FROM alpine:latest

RUN apk --no-cache add curl

WORKDIR /app

COPY --from=base /app/mock-app/_output/mock-app ./

EXPOSE 9000

# HEALTHCHECK --start-period=2s --interval=5s CMD curl -s --fail http://localhost:9000/status || exit 1

CMD [ "./mock-app" ]

