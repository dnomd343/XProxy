package config

var defaultConfig = `# default configure file for xproxy
log: debug

proxy:
  sniff: true
  redirect: true

network:
  dns: null
  ipv4: null
  ipv6: null
  bypass:
    - 169.254.0.0/16
    - 224.0.0.0/3
    - fc00::/7
    - fe80::/10
    - ff00::/8

update:
  cron: "0 0 4 * * *"
  url:
    geoip.dat: "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat"
    geosite.dat: "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat"
`
