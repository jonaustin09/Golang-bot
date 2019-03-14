#!/usr/bin/env bash

echo "Fetch dependencies..."
export GO111MODULE=on
go mod tidy

echo "Apply migrations..."
cd migrations
goose sqlite3 ../db.sqlite3 up
cd ..

echo "Remove old build..."
rm --f build

echo "Creating new build..."
go build -o build

echo "Run..."
./build