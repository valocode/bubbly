FROM golang:1.15.2-buster

WORKDIR /src

COPY go.sum .
COPY go.mod .

RUN go mod download

COPY . .
WORKDIR /src/storefront

CMD ["go", "run", "main.go"]

