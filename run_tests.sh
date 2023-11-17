#!/bin/bash

echo -e "\e[1mRunning unit and integration tests. Note: Some integration tests may take longer to execute, so please be patient.\e[0m"

for dir in */; do
    cd $dir
    echo "Running tests in $dir.."
    go test -v ./...
    cd ..
done
