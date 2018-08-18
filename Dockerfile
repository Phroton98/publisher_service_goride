FROM golang:alpine

RUN apk add --update tzdata bash wget curl git;
RUN mkdir -p /go/src/bin && curl https://glide.sh/get | sh
ADD . /go/src/app.goride
WORKDIR /go/src/app.goride
CMD glide update
CMD go run /go/src/app.goride/main.go