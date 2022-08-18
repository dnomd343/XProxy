package radvd

import (
    "XProxy/cmd/common"
    log "github.com/sirupsen/logrus"
    "strings"
)

func optionList(options map[string]string, intendNum int) string {
    var result string
    intend := strings.Repeat(" ", intendNum)
    for option, value := range options {
        result += intend + option + " " + value + ";\n"
    }
    return result
}

func loadPrefix(prefix string, options map[string]string) string {
    result := "    prefix " + prefix + " {\n"
    result += optionList(options, 8)
    return result + "    };\n"
}

func Load(options map[string]string, prefixes map[string]map[string]string) {
    radvdConfig := "interface eth0 {\n"
    radvdConfig += optionList(options, 4)
    for prefix, prefixOptions := range prefixes {
        radvdConfig += loadPrefix(prefix, prefixOptions)
    }
    radvdConfig += "};\n"
    log.Debugf("Radvd configure -> \n%s", radvdConfig)
    common.WriteFile("/etc/radvd.conf", radvdConfig, true)
}
