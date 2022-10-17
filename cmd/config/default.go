package config

import (
    "XProxy/cmd/common"
    "bytes"
    "encoding/json"
    "github.com/BurntSushi/toml"
    log "github.com/sirupsen/logrus"
    "gopkg.in/yaml.v3"
    "path"
)

var defaultConfig = map[string]interface{}{
    // TODO: add proxy bin option
    "proxy": map[string]string{
        //"core": "xray",
        "log": "warning",
    },
    "network": map[string]interface{}{
        "bypass": []string{
            "169.254.0.0/16",
            "224.0.0.0/3",
            "fc00::/7",
            "fe80::/10",
            "ff00::/8",
        },
    },
    "asset": map[string]interface{}{
        "update": map[string]interface{}{
            "cron": "0 5 6 * * *",
            "url": map[string]string{
                "geoip.dat":   "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat",
                "geosite.dat": "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat",
            },
        },
    },
}

func toJSON(config interface{}) string { // convert to JSON string
    jsonRaw, _ := json.MarshalIndent(config, "", "  ")
    return string(jsonRaw)
}

func toYAML(config interface{}) string { // convert to YAML string
    buf := new(bytes.Buffer)
    encoder := yaml.NewEncoder(buf)
    encoder.SetIndent(2) // with 2 space indent
    _ = encoder.Encode(config)
    return buf.String()
}

func toTOML(config interface{}) string { // convert to TOML string
    buf := new(bytes.Buffer)
    _ = toml.NewEncoder(buf).Encode(config)
    return buf.String()
}

func loadDefaultConfig(configFile string) {
    log.Infof("Load default configure -> %s", configFile)
    suffix := path.Ext(configFile)
    if suffix == ".json" {
        common.WriteFile(configFile, toJSON(defaultConfig), false) // JSON format
    } else if suffix == ".toml" {
        common.WriteFile(configFile, toTOML(defaultConfig), false) // TOML format
    } else {
        common.WriteFile(configFile, toYAML(defaultConfig), false) // YAML format
    }
}
