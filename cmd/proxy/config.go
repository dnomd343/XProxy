package proxy

import (
    "XProxy/cmd/common"
    log "github.com/sirupsen/logrus"
)

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

type logObject struct {
    Log struct {
        Loglevel string `json:"loglevel"`
        Access   string `json:"access"`
        Error    string `json:"error"`
    } `json:"log"`
}

type inboundsObject struct {
    Inbounds []interface{} `json:"inbounds"`
}

type sniffObject struct {
    Enabled         bool     `json:"enabled"`
    RouteOnly       bool     `json:"routeOnly"`
    DestOverride    []string `json:"destOverride"`
    DomainsExcluded []string `json:"domainsExcluded"`
}

type inboundObject struct {
    Tag            string      `json:"tag"`
    Port           int         `json:"port"`
    Protocol       string      `json:"protocol"`
    Settings       interface{} `json:"settings"`
    StreamSettings interface{} `json:"streamSettings"`
    Sniffing       sniffObject `json:"sniffing"`
}

func loadLogConfig(logLevel string, logDir string) string {
    if logLevel != "debug" && logLevel != "info" &&
        logLevel != "warning" && logLevel != "error" && logLevel != "none" {
        log.Warningf("Unknown log level -> %s", logLevel)
        logLevel = "warning" // using `warning` as default
    }
    logConfig := logObject{}
    logConfig.Log.Loglevel = logLevel
    logConfig.Log.Access = logDir + "/access.log"
    logConfig.Log.Error = logDir + "/error.log"
    return common.JsonEncode(logConfig)
}

func loadHttpConfig(tag string, port int, sniff sniffObject) interface{} {
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

func loadSocksConfig(tag string, port int, sniff sniffObject) interface{} {
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

func loadTProxyConfig(tag string, port int, sniff sniffObject) interface{} {
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
