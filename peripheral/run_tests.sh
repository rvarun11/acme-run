#!/bin/bash

echo -e "\e[1mWelcome to Peripheral Service Test Suite.\e[0m"
echo "The test suite only focuses on critical functionalities for key scenarios. It includes both unit and integration tests."
echo "Some integration tests may take longer; your patience is appreciated."
echo "Running tests..."

echo -e "\n"
go test -v ./...
echo -e "\n"

echo -e "\e[1mPeripheral Service Test Suite has completed."
