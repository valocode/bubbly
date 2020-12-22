FROM golang:1.15.2-buster AS builder

WORKDIR /src

COPY go.sum .
COPY go.mod .

RUN go mod download

COPY . .

# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS

# compile the go binary
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /bubbly

# compile the go integration tests
RUN go test ./integration -tags=integration,incluster -c -o ./integration/runtests

# target integration
FROM debian:buster-slim AS integration
COPY --from=builder /src/integration /test

WORKDIR /test

# some packages like certs are needed for our tests
RUN apt update -qq \
    && apt install -y -qq ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ENTRYPOINT [ "./runtests" ]

# target for running in cluster
FROM debian:buster-slim
COPY --from=builder /bubbly /bubbly

# Define GOTRACEBACK to mark this container as using the Go language runtime
# for `skaffold debug` (https://skaffold.dev/docs/workflows/debug/).
ENV GOTRACEBACK=single

ENTRYPOINT ["/bubbly"]
CMD ["--help"]

EXPOSE 8111
