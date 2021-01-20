FROM golang:1.15.2-buster AS builder

WORKDIR /src

COPY go.sum .
COPY go.mod .

RUN go mod download

COPY . .

# compile the go binary
RUN go build -o /bubbly

# target for running in cluster
FROM debian:buster-slim
COPY --from=builder /bubbly /bubbly

ENTRYPOINT ["/bubbly"]
CMD ["--help"]

EXPOSE 8111
