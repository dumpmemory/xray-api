package api

import (
	"os"
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
            "listen":"127.0.0.1",
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
    ],
    "dns": {
        "servers": [
            "223.5.5.5","119.29.29.29",
            "1.0.0.1","8.8.4.4",
            "1.1.1.1","8.8.8.8",
            "https+local://cloudflare-dns.com/dns-query",
            "localhost"
        ]
    },
    "routing": {
        "domainStrategy": "AsIs",
        "settings": {
            "rules": [
                {
                    "type": "field",
                    "inboundTag": ["api"],
                    "outboundTag": "api"
                },
                {
                    "type": "field",
                    "ip": ["0.0.0.0/8","10.0.0.0/8","100.64.0.0/10","127.0.0.0/8","169.254.0.0/16","172.16.0.0/12","192.0.0.0/24","192.0.2.0/24","192.168.0.0/16","198.18.0.0/15","198.51.100.0/24","203.0.113.0/24","::1/128","fc00::/7","fe80::/10"],
                    "outboundTag": "blocked"
                },
                {
                    "type": "field",
                    "domain": ["domain:epochtimes.com","domain:epochtimes.com.tw","domain:epochtimes.fr","domain:epochtimes.de","domain:epochtimes.jp","domain:epochtimes.ru","domain:epochtimes.co.il","domain:epochtimes.co.kr","domain:epochtimes-romania.com","domain:erabaru.net","domain:lagranepoca.com","domain:theepochtimes.com","domain:ntdtv.com","domain:ntd.tv","domain:ntdtv-dc.com","domain:ntdtv.com.tw","domain:minghui.org","domain:renminbao.com","domain:dafahao.com","domain:dongtaiwang.com","domain:falundafa.org","domain:wujieliulan.com","domain:ninecommentaries.com","domain:shenyun.com"],
                    "outboundTag": "blocked"
                },
                {
                    "type": "field",
                    "protocol": ["bittorrent"],
                    "outboundTag": "blocked"
                }
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
	f, _ := os.Create(path)
	tpl, _ := template.New("config").Parse(tplstr)
	err = tpl.Execute(f, xray)
	return
}
