#!/usr/bin/env bash

export GO111MODULE=on

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