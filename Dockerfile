FROM golang:latest

COPY . /usr/local/app

WORKDIR /usr/local/app

RUN go get ./...
