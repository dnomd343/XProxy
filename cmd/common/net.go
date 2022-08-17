package common

import (
    "net"
    "strings"
)

func isIP(ipAddr string, isCidr bool) bool {
    if !isCidr {
        return net.ParseIP(ipAddr) != nil
    }
    _, _, err := net.ParseCIDR(ipAddr)
    return err == nil
}

func IsIPv4(ipAddr string, isCidr bool) bool {
    return isIP(ipAddr, isCidr) && strings.Contains(ipAddr, ".")
}

func IsIPv6(ipAddr string, isCidr bool) bool {
    return isIP(ipAddr, isCidr) && strings.Contains(ipAddr, ":")
}
