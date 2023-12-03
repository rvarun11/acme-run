FROM golang:1.21.3-bookworm

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o challenge ./cmd

CMD ["sh", "-c", "./run_tests.sh && ./challenge"]
