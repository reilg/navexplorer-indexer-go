FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

RUN adduser -D -u 1001 -g '' appuser
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/indexerd ./cmd/indexerd
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/indexer-cli ./cmd/indexer-cli
RUN chmod u+x /go/bin/*

FROM scratch

WORKDIR /app

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/indexerd /app/indexerd
COPY --from=builder /go/bin/indexer-cli /app/indexer-cli

COPY .env.dist /app/.env

COPY ./config/mappings /app/mappings

ENTRYPOINT ["/app/indexerd"]