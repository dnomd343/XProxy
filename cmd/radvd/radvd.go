package radvd

import (
    "XProxy/cmd/common"
    log "github.com/sirupsen/logrus"
    "strings"
)

type Config struct {
    Enable bool              `yaml:"enable" json:"enable"`
    Client []string          `yaml:"client" json:"client"`
    Option map[string]string `yaml:"option" json:"option"`
    Route  struct {
        Cidr   string            `yaml:"cidr" json:"cidr"`
        Option map[string]string `yaml:"option" json:"option"`
    } `yaml:"route" json:"route"`
    Prefix struct {
        Cidr   string            `yaml:"cidr" json:"cidr"`
        Option map[string]string `yaml:"option" json:"option"`
    } `yaml:"prefix" json:"prefix"`
    DNSSL struct { // DNS Search List
        Suffix []string          `yaml:"suffix" json:"suffix"`
        Option map[string]string `yaml:"option" json:"option"`
    } `yaml:"dnssl" json:"dnssl"`
    RDNSS struct { // Recursive DNS Server
        IP     []string          `yaml:"ip" json:"ip"`
        Option map[string]string `yaml:"option" json:"option"`
    } `yaml:"rdnss" json:"rdnss"`
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

func loadClient(clients []string) string {
    if len(clients) == 0 { // without client settings
        return ""
    }
    ret := genSpace(4) + "clients {\n"
    for _, client := range clients {
        ret += genSpace(8) + client + ";\n"
    }
    return ret + genSpace(4) + "};\n"
}

func loadPrefix(prefix string, option map[string]string) string { // load radvd prefix configure
    if prefix == "" {                                             // without prefix settings
        return ""
    }
    header := genSpace(4) + "prefix " + prefix + " {\n"
    return header + loadOption(option, 8) + genSpace(4) + "};\n"
}

func loadRoute(cidr string, option map[string]string) string { // load radvd route configure
    if cidr == "" {                                            // without route settings
        return ""
    }
    header := genSpace(4) + "route " + cidr + " {\n"
    return header + loadOption(option, 8) + genSpace(4) + "};\n"
}

func loadRdnss(ip []string, option map[string]string) string {
    if len(ip) == 0 { // without rdnss settings
        return ""
    }
    header := genSpace(4) + "RDNSS " + strings.Join(ip, " ") + " {\n"
    return header + loadOption(option, 8) + genSpace(4) + "};\n"
}

func loadDnssl(suffix []string, option map[string]string) string {
    if len(suffix) == 0 { // without dnssl settings
        return ""
    }
    header := genSpace(4) + "DNSSL " + strings.Join(suffix, " ") + " {\n"
    return header + loadOption(option, 8) + genSpace(4) + "};\n"
}

func Load(Radvd *Config) {
    config := "interface eth0 {\n"
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
