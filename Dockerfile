FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=auto
RUN export GOFLAGS=-mod=vendor
RUN git clone https://github.com/MuhammadMahdiHusayniX/go-todolist.git go-todolist

RUN cd ./go-todolist && go get ./...
RUN cd ./go-todolist && go build

WORKDIR /build/go-todolist
EXPOSE 3000

ENTRYPOINT ["./go-todolist"]