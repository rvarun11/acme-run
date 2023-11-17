#!/bin/bash

echo -e "\e[1mTest Information:\e[0m The test suite focuses exclusively on critical functionalities essential for key scenarios."
echo -e "\e[1mRunning unit and integration tests. Note: Some integration tests may have longer execution times, so your patience is appreciated.\e[0m"
echo -e "\n"

for dir in */; do
    cd $dir
    echo -e "\e[1mRunning tests in $dir\e[0m"
    go test -v ./...
    cd ..
    echo -e "\n"
done
