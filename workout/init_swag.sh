#!/bin/bash

# Add the folders that have swagger annotations
swag fmt
swag init -g cmd/main.go internal/adapters/handler/http/http.go
