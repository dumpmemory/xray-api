package api

import (
	"encoding/json"
	"io"
	"text/template"
	"xray-api/config"
)

var tplstr = `{
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
            "port": {{.Grpc}},
            "protocol": "dokodemo-door",
            "settings": {"address": "127.0.0.1"}
        },
        {
            "tag": "vryusers",
            "port": {{.Port}},
            "protocol": "{{.Protocol}}",            
        {{if eq .Protocol "vmess"}}
            "settings": {"clients": []},
            {{if eq .Type "ws"}}
                "streamSettings": {
                    "network": "ws",
                    "security": "none",
                    "wsSettings": {"path": "{{.Path}}"},
                    "sockopt": {"mark": 0,"tcpFastOpen": true,"tproxy": "off"}
                },
            {{else if eq .Type "tcp"}}
                "streamSettings": {
                    "network": "tcp",
                    "security": {{if .Tls}}"tls"{{else}}"none"{{end}},
                    {{if .Tls}}
                    "tlsSettings": {
                        "certificates": [{
                            "certificateFile": "{{.CertificateFile}}",
                            "keyFile": "{{.KeyFile}}"
                        }]
                    },
                    {{end}}
                    "tcpSettings": {}
                },
            {{end}}
        {{else if eq .Protocol "shadowsocks"}}
            "settings":{
                "clients": [],
                "network": "{{.Network}}"
            },
		{{else if eq .Protocol "trojan"}}
			"settings": {"clients": []},
            "streamSettings": {
                "network": "tcp",
                "security": "tls",
                "tlsSettings": {
                    "alpn": ["http/1.1"],
                    "certificates": [{
						"certificateFile": "{{.CertificateFile}}",
						"keyFile": "{{.KeyFile}}"
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
    {{if gt (len .Dns) 0}},"dns": {"servers": {{js .Dns}}}{{end}},
    "routing": {
        "domainStrategy": "AsIs",
        "settings": {
            "rules": [
                {
                    "type": "field",
                    "inboundTag": ["api"],
                    "outboundTag": "api"
                }{{if gt (len .Block.Ips) 0}},{
                    "type": "field",
                    "ip": {{js .Block.Ips}},
                    "outboundTag": "blocked"
                }{{end}}{{if gt (len .Block.Domains) 0}},{
                    "type": "field",
                    "domain": {{js .Block.Domains}},
                    "outboundTag": "blocked"
                }{{end}}{{if .Block.BT}},
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

func JS(data interface{}) string {
	js, err := json.Marshal(data)
	if err != nil {
		return ""
	} else {
		return string(js)
	}
}
func WriteConfig(xray config.XRAY, wr io.Writer) (err error) {
	// if xray.Dns != nil {
	// 	var dns []string
	// 	for _, x := range xray["dns"].([]interface{}) {
	// 		dns = append(dns, `"`+x.(string)+`"`)
	// 	}
	// 	xray["dns"] = "[" + strings.Join(dns, ",") + "]"
	// }
	// if xray["blockIps"] != nil {
	// 	var blockIps []string
	// 	for _, x := range xray["blockIps"].([]interface{}) {
	// 		blockIps = append(blockIps, `"`+x.(string)+`"`)
	// 	}
	// 	xray["blockIps"] = "[" + strings.Join(blockIps, ",") + "]"
	// }
	// if xray["blockDomains"] != nil {
	// 	var blockDomains []string
	// 	for _, x := range xray["blockDomains"].([]interface{}) {
	// 		blockDomains = append(blockDomains, `"domain:`+x.(string)+`"`)
	// 	}
	// 	xray["blockDomains"] = "[" + strings.Join(blockDomains, ",") + "]"
	// }

	tpl, _ := template.New("config").Funcs(template.FuncMap{"js": JS}).Parse(tplstr)
	err = tpl.Execute(wr, xray)
	return
}
