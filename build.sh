#!/bin/sh
rm -rf release
mkdir release
cp xray-core release/xray-core -r
cp config.yaml release/config.yaml
cp config.yaml release/config.default.yaml
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o release/xray-api -ldflags='-s -w'
rm xray-api.zip
cd release
zip -r ../xray-api.zip *
cd ..