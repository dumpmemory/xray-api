#!/bin/sh
rm -rf release
mkdir release
cp xray-core release/xray-core -r
cp config.yaml release/config.yaml
cp config.yaml release/config.default.yaml
go build -o release/xray-api .