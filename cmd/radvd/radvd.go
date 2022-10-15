package radvd

import (
    "XProxy/cmd/common"
    log "github.com/sirupsen/logrus"
    "strings"
)

type Config struct {
    Log    int               `yaml:"log" json:"log" toml:"log"`
    Dev    string            `yaml:"dev" json:"dev" toml:"dev"`
    Enable bool              `yaml:"enable" json:"enable" toml:"enable"`
    Client []string          `yaml:"client" json:"client" toml:"client"`
    Option map[string]string `yaml:"option" json:"option" toml:"option"`
    Route  struct {
        Cidr   string            `yaml:"cidr" json:"cidr" toml:"cidr"`
        Option map[string]string `yaml:"option" json:"option" toml:"option"`
    } `yaml:"route" json:"route" toml:"route"`
    Prefix struct {
        Cidr   string            `yaml:"cidr" json:"cidr" toml:"cidr"`
        Option map[string]string `yaml:"option" json:"option" toml:"option"`
    } `yaml:"prefix" json:"prefix" toml:"prefix"`
    DNSSL struct { // DNS Search List
        Suffix []string          `yaml:"suffix" json:"suffix" toml:"suffix"`
        Option map[string]string `yaml:"option" json:"option" toml:"option"`
    } `yaml:"dnssl" json:"dnssl" toml:"dnssl"`
    RDNSS struct { // Recursive DNS Server
        IP     []string          `yaml:"ip" json:"ip" toml:"ip"`
        Option map[string]string `yaml:"option" json:"option" toml:"option"`
    } `yaml:"rdnss" json:"rdnss" toml:"rdnss"`
}

func genSpace(num int) string {
    return strings.Repeat(" ", num)
}

func loadOption(options map[string]string, intend int) string { // load options into radvd config format
    var ret string
    for option, value := range options {
        ret += genSpace(intend) + option + " " + value + ";\n"
    }
    return ret
}

func loadClient(clients []string) string { // load radvd client configure
    if len(clients) == 0 {
        return "" // without client settings
    }
    ret := genSpace(4) + "clients {\n"
    for _, client := range clients {
        ret += genSpace(8) + client + ";\n"
    }
    return ret + genSpace(4) + "};\n"
}

func loadPrefix(prefix string, option map[string]string) string { // load radvd prefix configure
    if prefix == "" {
        return "" // without prefix settings
    }
    header := genSpace(4) + "prefix " + prefix + " {\n"
    return header + loadOption(option, 8) + genSpace(4) + "};\n"
}

func loadRoute(cidr string, option map[string]string) string { // load radvd route configure
    if cidr == "" {
        return "" // without route settings
    }
    header := genSpace(4) + "route " + cidr + " {\n"
    return header + loadOption(option, 8) + genSpace(4) + "};\n"
}

func loadRdnss(ip []string, option map[string]string) string { // load radvd RDNSS configure
    if len(ip) == 0 {
        return "" // without rdnss settings
    }
    header := genSpace(4) + "RDNSS " + strings.Join(ip, " ") + " {\n"
    return header + loadOption(option, 8) + genSpace(4) + "};\n"
}

func loadDnssl(suffix []string, option map[string]string) string { // load radvd DNSSL configure
    if len(suffix) == 0 {
        return "" // without dnssl settings
    }
    header := genSpace(4) + "DNSSL " + strings.Join(suffix, " ") + " {\n"
    return header + loadOption(option, 8) + genSpace(4) + "};\n"
}

func Load(Radvd *Config) {
    config := "interface " + Radvd.Dev + " {\n"
    config += loadOption(Radvd.Option, 4)
    config += loadPrefix(Radvd.Prefix.Cidr, Radvd.Prefix.Option)
    config += loadRoute(Radvd.Route.Cidr, Radvd.Route.Option)
    config += loadClient(Radvd.Client)
    config += loadRdnss(Radvd.RDNSS.IP, Radvd.RDNSS.Option)
    config += loadDnssl(Radvd.DNSSL.Suffix, Radvd.DNSSL.Option)
    config += "};\n"
    log.Debugf("Radvd configure -> \n%s", config)
    common.WriteFile("/etc/radvd.conf", config, true)
}
