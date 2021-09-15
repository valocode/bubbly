# Stage 1 build
FROM node:16.8.0-buster as ui

WORKDIR /work

COPY ui/package*.json ./

RUN npm install

COPY ui .

RUN npm run build

FROM golang:1.17-buster AS builder

WORKDIR $GOPATH/src/github.com/valocode/bubbly

COPY go.sum .
COPY go.mod .

RUN go mod download \
    && go mod verify

COPY . .

# Copy the built UI over
COPY --from=ui /work/build ui/build

RUN go build -tags ui -o /go/bin/bubbly

# step 2 deploy
FROM gcr.io/distroless/base-debian10

# Copy our static executable.
COPY --from=builder /go/bin/bubbly /bubbly

# Use an unprivileged user.
USER nonroot:nonroot

EXPOSE 8111
ENTRYPOINT ["/bubbly"]
CMD ["server"]


