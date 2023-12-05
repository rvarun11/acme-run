#!/bin/bash

echo -e "\e[1mWelcome to ACME Run Test Suite by Liuyin, Samkith & Varun.\e[0m"
echo "The test suite only focuses on critical functionalities for key scenarios."
echo "Some tests may take longer; your patience is appreciated."
echo "Running tests..."

for dir in */; do
    cd $dir
    echo -e "\e[1mRunning tests in $dir\e[0m"
    go test -v ./...
    cd ..
done

echo -e "\e[1mAll tests completed.\e[0m"
