FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=on
RUN export GOFLAGS=-mod=vendor
RUN git clone https://github.com/MuhammadMahdiHusayniX/go-todolist.git

RUN cd /go-todolist && go get ./...
RUN go build

EXPOSE 3000

ENTRYPOINT ['/build/go-todolist/go-todolist']