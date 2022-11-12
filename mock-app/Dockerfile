# syntax=docker/dockerfile:1

FROM golang:1.19.3-alpine AS base

RUN apk --no-cache add git
RUN apk --no-cache add bash

WORKDIR /app/mock-app

COPY ./mock-app/go.* ./

RUN go mod download

WORKDIR /app

COPY . .

RUN chmod +x ./mock-app/build.sh

RUN ./mock-app/build.sh

FROM alpine:latest

WORKDIR /app

COPY --from=base /app/_output/mock-app ./

EXPOSE 9000

CMD [ "./mock-app" ]
