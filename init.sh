wget https://github.com/XTLS/Xray-core/releases/download/v1.4.2/Xray-linux-64.zip
unzip -o Xray-linux-64.zip -d xray-core
rm Xray-linux-64.zip
go mod download
go get -u