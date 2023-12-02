#!/bin/bash

# Add the folders that have swagger annotations
go install github.com/swaggo/swag/cmd/swag@latest
swag fmt
swag init -g cmd/main.go internal/adapters/handler/http/http.go
