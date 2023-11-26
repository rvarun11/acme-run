#!/bin/bash

# Add the folders that have swagger annotatons
swag fmt
swag init -g cmd/main.go internal/adapters/handler/http/http.go
