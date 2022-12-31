FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/server
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go mod download
RUN go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/server .

RUN chmod +x /go/bin/server

FROM scratch

COPY --from=builder /go/bin/server /go/bin/server

EXPOSE 8080
ENTRYPOINT ["/go/bin/server"]
