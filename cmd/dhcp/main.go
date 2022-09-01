package dhcp

type dhcpConfig struct {
    Enable    bool   `yaml:"enable" json:"enable"`
    Configure string `yaml:"config" json:"config"`
}

type Config struct {
    IPv4 dhcpConfig `yaml:"ipv4" json:"ipv4"`
    IPv6 dhcpConfig `yaml:"ipv6" json:"ipv6"`
}
