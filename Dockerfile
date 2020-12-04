FROM golang:1.15.2-buster AS builder

WORKDIR /src

COPY go.sum .
COPY go.mod .

RUN go mod download

COPY . .

# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /bubbly

FROM golang:1.15.2-buster AS tester
COPY --from=builder /src /src

WORKDIR /src

FROM debian:buster-slim
COPY --from=builder /bubbly /bubbly

# Define GOTRACEBACK to mark this container as using the Go language runtime
# for `skaffold debug` (https://skaffold.dev/docs/workflows/debug/).
ENV GOTRACEBACK=single

ENTRYPOINT ["/bubbly"]
CMD ["--help"]

EXPOSE 8111
