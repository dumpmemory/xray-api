#!/bin/sh
wget https://github.com/v2fly/xray-core/releases/download/v4.34.0/xray-linux-64.zip
unzip -o xray-linux-64.zip -d main/xray-core
rm xray-linux-64.zip
go get github.com/xtls/xray-core v1.3.0
go get google.golang.org/grpc v1.35.0
go get github.com/pires/go-proxyproto
go get github.com/seiflotfy/cuckoofilter
go get github.com/akkuman/parseConfig
go get github.com/gin-gonic/gin

git clone git@github.com:zcmimi/xray-api-release.git expt