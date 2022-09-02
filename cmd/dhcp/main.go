package dhcp

import (
    "XProxy/cmd/common"
    log "github.com/sirupsen/logrus"
    "path"
)

var WorkDir = "/etc/dhcp"

type dhcpConfig struct {
    Enable    bool   `yaml:"enable" json:"enable"`
    Configure string `yaml:"config" json:"config"`
}

type Config struct {
    IPv4 dhcpConfig `yaml:"ipv4" json:"ipv4"`
    IPv6 dhcpConfig `yaml:"ipv6" json:"ipv6"`
}

func Load(config *Config) {
    if config.IPv4.Enable {
        log.Infof("Load DHCPv4 configure")
        common.WriteFile(path.Join(WorkDir, "dhcp4.conf"), config.IPv4.Configure, true)
    }
    if config.IPv6.Enable {
        log.Infof("Load DHCPv6 configure")
        common.WriteFile(path.Join(WorkDir, "dhcp6.conf"), config.IPv6.Configure, true)
    }
}
