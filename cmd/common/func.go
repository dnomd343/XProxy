package common

import (
    "encoding/json"
    log "github.com/sirupsen/logrus"
    "net"
    "os/exec"
    "strings"
    "syscall"
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

func JsonEncode(raw interface{}) string {
    jsonOutput, _ := json.MarshalIndent(raw, "", "  ") // json encode
    return string(jsonOutput)
}

func RunCommand(command ...string) (int, string) {
    log.Debugf("Running system command -> %v", command)
    process := exec.Command(command[0], command[1:]...)
    output, _ := process.CombinedOutput()
    log.Debugf("Command %v -> \n%s", command, string(output))
    code := process.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
    if code != 0 {
        log.Warningf("Command %v return code %d", command, code)
    }
    return code, string(output)
}
