FROM golang
COPY . /go/src/instore_server
WORKDIR /go/src/instore_server
RUN go get ./server
ENTRYPOINT go run server/server.go
EXPOSE 50051
