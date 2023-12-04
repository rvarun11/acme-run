#!/bin/bash

echo -e "\e[1mWelcome to Challenge Manager Test Suite.\e[0m"
echo "The test suite only focuses on critical functionalities for key scenarios."
echo "Some tests may take longer; your patience is appreciated."
echo "Running tests..."

go test -v ./...

echo -e "\e[1mChallenge Manager Test Suite has completed."
