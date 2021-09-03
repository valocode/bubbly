# Stage 1 build
FROM golang:1.16-buster AS builder


WORKDIR $GOPATH/src/github.com/valocode/bubbly

COPY go.sum .
COPY go.mod .

RUN go mod download \
    && go mod verify

COPY . .

RUN go build -o /go/bin/bubbly

# step 2 deploy
FROM gcr.io/distroless/base-debian10

# Copy our static executable.
COPY --from=builder /go/bin/bubbly go/bin/bubbly

# Use an unprivileged user.
USER nonroot:nonroot

ENTRYPOINT ["go/bin/bubbly"]
# 4223 NATS service
# 8111 bubbly agent
# 8222 NATS HTTP
EXPOSE 4223 8111 8222
