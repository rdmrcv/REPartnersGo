FROM golang:1.26.0-alpine3.23

RUN make

FROM alpine:3.23.0

