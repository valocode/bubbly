FROM golang:1.16.1-buster AS builder

WORKDIR /src

COPY go.sum .
COPY go.mod .

RUN go mod download

RUN apt-get update && apt-get install ca-certificates -y

COPY . .

RUN go get -u github.com/swaggo/swag/cmd/swag
RUN swag init

# compile the go binary
RUN go build -o /bubbly

# target for running in cluster
FROM debian:buster-slim
COPY --from=builder /bubbly /bubbly
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/bubbly"]
CMD ["--help"]

# 4223 NATS service
# 8111 bubbly agent
# 8222 NATS HTTP
EXPOSE 4223 8111 8222