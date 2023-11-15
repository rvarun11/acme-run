FROM golang:1.21.3-bookworm

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o workout ./cmd

EXPOSE 8005

CMD ["./trail"]
