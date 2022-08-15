package main

import (
    "fmt"
    "net"
    "strconv"
    "strings"
)

func isIP(ipAddr string, isRange bool, allowEmpty bool, ipLength int, ipFlag string) bool {
    var address string
    if allowEmpty && ipAddr == "" { // empty case
        return true
    }
    if isRange {
        temp := strings.Split(ipAddr, "/")
        if len(temp) != 2 { // not {IP_ADDRESS}/{LENGTH} format
            return false
        }
        length, err := strconv.Atoi(temp[1])
        if err != nil { // range length not a integer
            return false
        }
        if length < 0 || length > ipLength { // length should between 0 ~ ipLength
            return false
        }
        address = temp[0]
    } else {
        address = ipAddr
    }
    ip := net.ParseIP(address) // try to convert ip
    return ip != nil && strings.Contains(address, ipFlag)
}

func isIPv4(ipAddr string, isRange bool, allowEmpty bool) bool {
    return isIP(ipAddr, isRange, allowEmpty, 32, ".")
}

func isIPv6(ipAddr string, isRange bool, allowEmpty bool) bool {
    return isIP(ipAddr, isRange, allowEmpty, 128, ":")
}

func main() {
    fmt.Println("XProxy")

    fmt.Println(isIPv4("1.1.1.1", false, false))
    fmt.Println(isIPv4("1.1.1.1/24", true, false))
    fmt.Println(isIPv4("", true, true))
    fmt.Println(isIPv4("::1", true, true))
    fmt.Println(isIPv6("::1", true, true))
    fmt.Println(isIPv6("::1/32", true, true))
    fmt.Println(isIPv6("::1", true, true))
}
