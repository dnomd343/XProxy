package proxy

var dnsConfig = `{
  "dns": {
    "servers": [
      "localhost"
    ]
  }
}`

var routeConfig = `{
  "routing": {
    "domainStrategy": "AsIs",
    "rules": [
      {
        "type": "field",
        "network": "tcp,udp",
        "outboundTag": "node"
      }
    ]
  }
}`

var outboundsConfig = `{
  "outbounds": [
    {
      "tag": "node",
      "protocol": "freedom",
      "settings": {}
    }
  ]
}`

func httpConfig(tag string, port int, sniff sniffObject) interface{} {
	type empty struct{}
	return inboundObject{
		Tag:            tag,
		Port:           port,
		Protocol:       "http",
		Settings:       empty{},
		StreamSettings: empty{},
		Sniffing:       sniff,
	}
}

func socksConfig(tag string, port int, sniff sniffObject) interface{} {
	type empty struct{}
	type socksObject struct {
		UDP bool `json:"udp"`
	}
	return inboundObject{
		Tag:            tag,
		Port:           port,
		Protocol:       "socks",
		Settings:       socksObject{UDP: true},
		StreamSettings: empty{},
		Sniffing:       sniff,
	}
}

func tproxyConfig(tag string, port int, sniff sniffObject) interface{} {
	type tproxyObject struct {
		Network        string `json:"network"`
		FollowRedirect bool   `json:"followRedirect"`
	}
	type tproxyStreamObject struct {
		Sockopt struct {
			Tproxy string `json:"tproxy"`
		} `json:"sockopt"`
	}
	tproxyStream := tproxyStreamObject{}
	tproxyStream.Sockopt.Tproxy = "tproxy"
	return inboundObject{
		Tag:      tag,
		Port:     port,
		Protocol: "dokodemo-door",
		Settings: tproxyObject{
			Network:        "tcp,udp",
			FollowRedirect: true,
		},
		StreamSettings: tproxyStream,
		Sniffing:       sniff,
	}
}
