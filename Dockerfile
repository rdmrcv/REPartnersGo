FROM golang:1.26.0-alpine3.23 AS builder

WORKDIR /tmp/build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -a -o out/server ./app

FROM alpine:3.23.0

COPY --from=builder /tmp/build/out/server /var/lib/solver/server
RUN chmod +x /var/lib/solver/server

ENTRYPOINT ["/var/lib/solver/server"]