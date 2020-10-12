FROM golang:1.15.2-alpine

WORKDIR /go/src/github.com/leviharrison/fios-exporter
COPY . .

RUN go get -d -v ./...
RUN go build cmd/exporter/exporter.go

ENTRYPOINT ["./exporter"]