FROM golang:1.25.3 AS build

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
RUN make build

FROM alpine:3.22.0 AS base
RUN apk add --no-cache rsync
COPY --from=build /src/fetcharr /fetcharr

ENTRYPOINT ["/fetcharr"]