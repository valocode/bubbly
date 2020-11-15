FROM golang:1.15.2-buster

WORKDIR /src

COPY go.sum .
COPY go.mod .

RUN go mod download

COPY . .
COPY docker/integration/test.bash .

CMD ["./test.bash"]

