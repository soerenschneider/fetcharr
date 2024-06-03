FROM golang:1.22.3 AS build

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
RUN make build

FROM alpine:3.19.1 AS base
RUN apk add --no-cache rsync
COPY --from=build /src/fetcharr /fetcharr

ENTRYPOINT ["/fetcharr"]