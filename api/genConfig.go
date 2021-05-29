package api

import (
	"os"
	"strings"
	"text/template"
)

const tplstr = `{
    "log": {
        "loglevel": "error"
    },
    "stats": {},
    "api": {
        "services": [
            "HandlerService",
            "LoggerService",
            "StatsService"
        ],
        "tag": "api"
    },
    "policy": {
        "levels": {
            "0": {"statsUserUplink": true,"statsUserDownlink": true}
        },
        "system": {
            "statsInboundUplink": true,"statsInboundDownlink": true}
    },
    "inbounds": [
        {
            "tag": "api",
            "listen":"127.0.0.1",
            "port": {{.grpcPort}},
            "protocol": "dokodemo-door",
            "settings": {"address": "127.0.0.1"}
        },
        {
            "tag": "vryusers",
            "port": {{.port}},
            "protocol": "{{.protocol}}",
            
        {{if eq .protocol "vmess"}}
            "settings": {"clients": []},
            {{if eq .type "ws"}}
                "streamSettings": {
                    "network": "ws",
                    "security": "none",
                    "wsSettings": {"path": "{{.path}}"},
                    "sockopt": {"mark": 0,"tcpFastOpen": true,"tproxy": "off"}
                },
            {{else if eq .type "tcp"}}
                "streamSettings": {
                    "network": "tcp",
                    "security": {{if .tls}}"tls"{{else}}"none"{{end}},
                    {{if .tls}}
                    "tlsSettings": {
                        "certificates": [{
                            "certificateFile": "{{.certificateFile}}",
                            "keyFile": "{{.keyFile}}"
                        }]
                    },
                    {{end}}
                    "tcpSettings": {}
                },
            {{end}}
        {{else if eq .protocol "shadowsocks"}}
            "settings":{
                "clients": [],
                "network": "{{.network}}"
            },
		{{else if eq .protocol "trojan"}}
			"settings": {"clients": []},
            "streamSettings": {
                "network": "tcp",
                "security": "tls",
                "tlsSettings": {
                    "alpn": ["http/1.1"],
                    "certificates": [{
						"certificateFile": "{{.certificateFile}}",
						"keyFile": "{{.keyFile}}"
					}]
                }
            }
        {{end}}
		
            "sniffing": {
                "enabled": true,
                "destOverride": ["http","tls"]
            },
            "allocate": {
                "strategy": "always",
                "refresh": 5,
                "concurrency": 3
            }
        }
    ],
    "outbounds": [
        {
            "protocol": "freedom",
            "settings": {
                "domainStrategy": "AsIs"
            },
            "tag": "direct"
        },
        {
            "protocol": "blackhole",
            "settings": {},
            "tag": "blocked"
        }
    ]
    {{if .dns}},"dns": {"servers": {{.dns}}}{{end}},
    "routing": {
        "domainStrategy": "AsIs",
        "settings": {
            "rules": [
                {
                    "type": "field",
                    "inboundTag": ["api"],
                    "outboundTag": "api"
                }{{if .blockDomains}},{
                    "type": "field",
                    "ip": {{.blockIps}},
                    "outboundTag": "blocked"
                }{{end}}{{if .blockDomains}},{
                    "type": "field",
                    "domain": {{.blockDomains}},
                    "outboundTag": "blocked"
                }{{end}}{{if .blockBt}},
                {
                    "type": "field",
                    "protocol": ["bittorrent"],
                    "outboundTag": "blocked"
                }{{end}}
            ]
        },
        "strategy": "rules"
    },
    "transport": {
        "kcpSettings": {
            "uplinkCapacity": 100,
            "downlinkCapacity": 100,
            "congestion": true
        }
    }
}
`

func WriteConfig(xray map[interface{}]interface{}, path string) (err error) {
	if xray["dns"] != nil {
		var dns []string
		for _, x := range xray["dns"].([]interface{}) {
			dns = append(dns, `"`+x.(string)+`"`)
		}
		xray["dns"] = "[" + strings.Join(dns, ",") + "]"
	}
	if xray["blockIps"] != nil {
		var blockIps []string
		for _, x := range xray["blockIps"].([]interface{}) {
			blockIps = append(blockIps, `"`+x.(string)+`"`)
		}
		xray["blockIps"] = "[" + strings.Join(blockIps, ",") + "]"
	}
	if xray["blockDomains"] != nil {
		var blockDomains []string
		for _, x := range xray["blockDomains"].([]interface{}) {
			blockDomains = append(blockDomains, `"domain:`+x.(string)+`"`)
		}
		xray["blockDomains"] = "[" + strings.Join(blockDomains, ",") + "]"
	}

	f, _ := os.Create(path)
	tpl, _ := template.New("config").Parse(tplstr)
	err = tpl.Execute(f, xray)
	return
}
