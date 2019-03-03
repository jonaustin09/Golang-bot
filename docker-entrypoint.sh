#!/bin/bash


cd /go/src/money/
echo "building"
go build -o build

# Start bot
echo "Starting bot"
./build
