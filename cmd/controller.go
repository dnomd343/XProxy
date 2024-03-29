package main

import (
    "XProxy/cmd/asset"
    "XProxy/cmd/common"
    "XProxy/cmd/config"
    "XProxy/cmd/dhcp"
    "XProxy/cmd/network"
    "XProxy/cmd/process"
    "XProxy/cmd/proxy"
    "XProxy/cmd/radvd"
    log "github.com/sirupsen/logrus"
    "os"
    "os/signal"
    "path"
    "strconv"
    "syscall"
    "time"
)

func runProcess(env []string, command ...string) {
    sub := process.New(command...)
    sub.Run(true, env)
    sub.Daemon()
    subProcess = append(subProcess, sub)
}

func blockWait() {
    sigExit := make(chan os.Signal, 1)
    signal.Notify(sigExit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM) // wait until get exit signal
    <-sigExit
}

func loadRadvd(settings *config.Config) {
    if settings.Radvd.Enable {
        radvd.Load(&settings.Radvd)
    } else {
        log.Infof("Skip loading radvd")
    }
}

func loadDhcp(settings *config.Config) {
    common.CreateFolder(dhcp.WorkDir)
    if settings.DHCP.IPv4.Enable || settings.DHCP.IPv6.Enable {
        common.CreateFolder(path.Join(exposeDir, "dhcp"))
        dhcp.Load(&settings.DHCP)
    }
}

func loadAsset(settings *config.Config) {
    if settings.Asset.Disable {
        log.Infof("Skip loading asset")
    } else {
        asset.Load(assetFile, assetDir)
        asset.AutoUpdate(&settings.Asset, assetDir)
    }
}

func loadNetwork(settings *config.Config) {
    settings.IPv4.RouteTable = v4RouteTable
    settings.IPv4.TProxyPort = v4TProxyPort
    settings.IPv6.RouteTable = v6RouteTable
    settings.IPv6.TProxyPort = v6TProxyPort
    network.Load(settings.DNS, settings.Dev, &settings.IPv4, &settings.IPv6)
}

func loadProxy(settings *config.Config) {
    if proxyBin != "" {
        settings.Proxy.Bin = proxyBin // setting proxy bin from env
    }
    settings.Proxy.V4TProxyPort = v4TProxyPort
    settings.Proxy.V6TProxyPort = v6TProxyPort
    proxy.Load(configDir, exposeDir, &settings.Proxy)
}

func runProxy(settings *config.Config) {
    assetEnv := []string{
        "XRAY_LOCATION_ASSET=" + assetDir,  // xray asset folder
        "V2RAY_LOCATION_ASSET=" + assetDir, // v2ray / sagray asset folder
    }
    runProcess(assetEnv, settings.Proxy.Bin, "run", "-confdir", configDir)
}

func runRadvd(settings *config.Config) {
    if settings.Radvd.Enable {
        radvdCmd := []string{"radvd", "--nodaemon"}
        if settings.Radvd.Log > 0 { // with log option
            radvdCmd = append(radvdCmd, "--logmethod", "logfile")
            radvdCmd = append(radvdCmd, "--logfile", path.Join(exposeDir, "log/radvd.log"))
            radvdCmd = append(radvdCmd, "--debug", strconv.Itoa(settings.Radvd.Log))
        }
        runProcess(nil, radvdCmd...)
    } else {
        log.Infof("Skip running radvd")
    }
}

func runDhcp(settings *config.Config) {
    leaseDir := path.Join(exposeDir, "dhcp")
    if settings.DHCP.IPv4.Enable {
        v4Leases := path.Join(leaseDir, "dhcp4.leases")
        v4Config := path.Join(dhcp.WorkDir, "dhcp4.conf")
        if !common.IsFileExist(v4Leases) {
            common.WriteFile(v4Leases, "", true)
        }
        runProcess(nil, "dhcpd", "-4", "-f", "-cf", v4Config, "-lf", v4Leases)
        time.Sleep(time.Second) // wait 1s for avoid cluttered output
    } else {
        log.Infof("Skip running DHCPv4")
    }
    if settings.DHCP.IPv6.Enable {
        v6Leases := path.Join(leaseDir, "dhcp6.leases")
        v6Config := path.Join(dhcp.WorkDir, "dhcp6.conf")
        if !common.IsFileExist(v6Leases) {
            common.WriteFile(v6Leases, "", true)
        }
        runProcess(nil, "dhcpd", "-6", "-f", "-cf", v6Config, "-lf", v6Leases)
        time.Sleep(time.Second) // wait 1s for avoid cluttered output
    } else {
        log.Infof("Skip running DHCPv6")
    }
}
