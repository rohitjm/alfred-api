# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY /actions ./actions
COPY /models ./models
COPY main.go ./
COPY country_codes.json ./

RUN go build -o /alfred-api

EXPOSE 10000

CMD [ "/alfred-api" ]
