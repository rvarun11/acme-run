FROM golang:1.21.3-bookworm

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o hrm ./cmd

# TODO: See if you need to expose the port or not
EXPOSE 8004

CMD ["./peripheral"]
