FROM golang:1.21.4 AS build

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
RUN make build

FROM alpine:3.19.0 AS base
RUN apk add --no-cache rsync
COPY --from=build /src/fetcharr /fetcharr

ENTRYPOINT ["/fetcharr"]