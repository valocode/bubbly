ARG GO_VERSION=1.15
FROM golang:${GO_VERSION}-alpine AS build-env

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GOOS=linux go build -o bubbly .

FROM alpine:latest

WORKDIR /dist

COPY --from=build-env /build/bubbly .


EXPOSE 8080

ENTRYPOINT [ "/dist/bubbly" ]
