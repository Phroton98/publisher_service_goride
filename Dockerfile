FROM golang:alpine

RUN apk add --update tzdata \
    bash wget curl git;

RUN mkdir -p $$GOPATH/bin && \
    curl https://glide.sh/get | sh