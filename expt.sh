#!/bin/sh
rm -rf expt/*
mkdir expt
mkdir expt/xray-core
go build -o expt/xray-api xray-api.go
cp config.yaml expt/config.yaml
cp config.yaml expt/config.default.yaml
cp install.sh expt/
cp readme.md expt/

cp xray-core/xray expt/xray-core/
cp xray-core/*.dat expt/xray-core/