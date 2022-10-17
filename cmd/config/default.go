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

var defaultConfig = `# default configure file for xproxy
proxy:
  core: xray
  log: warning

network:
  bypass:
    - 169.254.0.0/16
    - 224.0.0.0/3
    - fc00::/7
    - fe80::/10
    - ff00::/8

asset:
  update:
    cron: "0 5 6 * * *"
    url:
      geoip.dat: https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat
      geosite.dat: https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat
`

func toJSON(yamlConfig string) string { // YAML -> JSON
    var config interface{}
    if err := yaml.Unmarshal([]byte(yamlConfig), &config); err != nil {
        log.Panicf("Default config error -> %v", err)
    }
    jsonRaw, _ := json.Marshal(config)
    return string(jsonRaw)
}

func toTOML(yamlConfig string) string { // YAML -> TOML
    var config interface{}
    if err := yaml.Unmarshal([]byte(yamlConfig), &config); err != nil {
        log.Panicf("Default config error -> %v", err)
    }
    buf := new(bytes.Buffer)
    _ = toml.NewEncoder(buf).Encode(config)
    return buf.String()
}

func loadDefaultConfig(configFile string) {
    log.Infof("Load default configure -> %s", configFile)
    suffix := path.Ext(configFile)
    if suffix == ".json" {
        defaultConfig = toJSON(defaultConfig)
    } else if suffix == ".toml" {
        defaultConfig = toTOML(defaultConfig)
    }
    common.WriteFile(configFile, defaultConfig, false)
}
