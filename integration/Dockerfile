FROM golang:1.15.2-buster AS builder

WORKDIR /src

COPY go.sum .
COPY go.mod .

RUN go mod download

COPY . .

# compile the go integration tests
RUN go test ./integration -tags=integration -c -o ./integration/integration.test

# target integration
FROM debian:buster-slim AS integration
COPY --from=builder /src/integration /test

WORKDIR /test

# some packages like certs are needed for our tests
RUN apt update -qq \
    && apt install -y -qq ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ENTRYPOINT [ "./integration.test" ]
