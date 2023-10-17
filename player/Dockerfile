FROM golang:1.21.3-bookworm

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o player ./cmd

EXPOSE 8000

CMD ["./player"]
