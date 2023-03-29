# 使用 XProxy 绕过校园网认证登录

部分校园网在登录认证时需要 DNS 解析，因而在防火墙上允许 `TCP/53` 或 `UDP/53` 端口通行，借助这个漏洞，可将内网流量用 XProxy 代理并转发到公网服务器上，实现免认证、无限速的上网。

以下为一般情况下的网络拓扑：

![Network](./img/campus_network.png)

为了方便讲解，我们假设以下典型情况：

+ 校园网交换机无 IPv6 支持，同时存在 QoS；

+ 无认证时允许 53 端口通行，ICMP 流量无法通过；

+ 使用三台公网服务器负载均衡，其 53 端口上运行有代理服务；

+ 三台服务器只有一台支持 IPv4 与 IPv6 双栈，其余只支持 IPv4；

## 代理协议

从部署成本与便捷性方面考虑，socks 类代理是最合适的工具：无需修改服务器网卡路由表等配置，方便多级负载均衡，软件只在用户态运行，实测速度也相对 `IPSec` 、`L2TP` 等协议更有优势；但 socks 代理只接收 TCP 与 UDP 流量，ICMP 流量无法被直接代理（例如 PING 命令），不过大多数情况下我们不会用到公网 ICMP 流量，如果确实需要也可以曲线救国给它补上。

在选定代理类型后，我们需要考虑具体的传输方式，由于存在 QoS 问题，这里应该倾向于选择基于 TCP 的代理方式，同时为了避免校园网的流量审查，我们应该将流量加密传输。考虑到软路由性能一般较差，而自建的代理服务器无需考虑协议兼容性问题，这里更建议选择基于 XTLS 的传输方式，它避开了对 TLS 流量的二次加密，可以显著降低代理 https 流量时的性能开销，提升性能上限；至于延迟方面的问题，如果选择 `gRPC` 等协议，虽然有 0-rtt 握手的延迟优势，但这种场景下延迟一般不高（甚至服务器可以直接部署在校内），用微弱的延迟优势换取性能开销并不值得，且前者也可以开启 mux 多路复用来优化延迟。

既然我们已经选择 XTLS 方式，那使用轻量的无加密类型（在加密的 XTLS 隧道里传输）是当前网络的最优解，譬如 VLESS 或者 Trojan 协议，下面将用 VLESS + XTLS 代理进行配置演示；当然，具体的选择还是取决于您的实际应用场景，只要按需调整 XProxy 的配置文件即可。

## 初始化配置

> 分配 `192.168.2.0/24` 和 `fc00::/64` 给内网使用。

路由器 WAN 口接入学校交换机，构建一个 NAT 转换，代理流量在路由器转发后送到公网服务器的 53 端口上；假设内网中路由器地址为 `192.168.2.1` ，配置虚拟网关 IPv4 地址为 `192.168.2.2` ，IPv6 地址为 `fc00::2` ；在网关中，无论 IPv4 还是 IPv6 流量都会被透明代理，由于校园网无 IPv6 支持，数据被封装后只通过 IPv4 网络发送，代理服务器接收以后再将其解开，对于 IPv6 流量，这里相当于一个 `6to4` 隧道。

```bash
# 宿主机网卡假定为 eth0
$ ip link set eth0 promisc on
$ modprobe ip6table_filter
$ docker network create -d macvlan \
  --subnet=192.168.2.0/24 \  # 此处指定的参数为容器的默认网络配置
  --gateway=192.168.2.1 \
  --subnet=fc00::/64 \
  --gateway=fc00::1 \
  --ipv6 -o parent=eth0 macvlan
```

我们将配置文件保存在 `/etc/scutweb` 目录下，使用以下命令开启 XProxy 服务：

```bash
docker run --restart always \
  --privileged --network macvlan -dt \
  --name scutweb --hostname scutweb \
  --volume /etc/scutweb/:/xproxy/ \
  --volume /etc/timezone:/etc/timezone:ro \
  --volume /etc/localtime:/etc/localtime:ro \
  dnomd343/xproxy:latest
```

## 参数配置

我们将三台服务器分别称为 `nodeA` ，`nodeB` 与 `nodeC` ，其中只有 `nodeC` 支持IPv6网络；此外，我们在内网分别暴露 3 个 socks5 端口，分别用于检测服务器的可用性。

由于校园网无 IPv6 支持，这里 IPv6 上游网关可以不填写；虚拟网关对内网发布 RA 通告，让内网设备使用 SLAAC 配置网络地址，同时将其作为 IPv6 网关；此外，如果路由器开启了 DHCP 服务，需要将默认网关改为 `192.168.2.2` ，也可以启用 XProxy 自带的 DHCPv4 服务。

最后，由于我们代理全部流量，无需根据域名或者 IP 进行任何分流，因此路由资源自动更新部分可以省略。

修改 `xproxy.yml` ，写入以下配置：

```yaml
proxy:
  log: warning
  socks:
    nodeA: 1081
    nodeB: 1082
    nodeC: 1083

network:
  dev: eth0
  dns:
    - 192.168.2.1
  ipv4:
    gateway: 192.168.2.1
    address: 192.168.2.2/24
  ipv6:
    gateway: null
    address: fc00::2/64
  bypass:
    - 169.254.0.0/16
    - 224.0.0.0/3
    - fc00::/7
    - fe80::/10
    - ff00::/8

radvd:
  log: 5
  dev: eth0
  enable: true
  option:
    AdvSendAdvert: on
  prefix:
    cidr: fc00::/64

custom:
  pre:
    - "iptables -t nat -N FAKE_PING"
    - "iptables -t nat -A FAKE_PING -j DNAT --to-destination 192.168.2.2"
    - "iptables -t nat -A PREROUTING -i eth0 -p icmp -j FAKE_PING"
    - "ip6tables -t nat -N FAKE_PING"
    - "ip6tables -t nat -A FAKE_PING -j DNAT --to-destination fc00::2"
    - "ip6tables -t nat -A PREROUTING -i eth0 -p icmp -j FAKE_PING"
```

在开始代理前，我们使用 `custom` 注入了一段脚本配置：由于这里我们只代理 TCP 与 UDP 流量，ICMP 数据包不走代理，内网设备 ping 外网时会一直无响应，加入这段脚本可以创建一个 NAT，假冒远程主机返回成功回复，但实际上 ICMP 数据包并未实际到达，效果上表现为 ping 成功且延迟为内网访问时间。

> 这段脚本并无实质作用，仅用于演示 `custom` 功能。

## 代理配置

接下来，我们应该配置出站代理，修改 `config/outbounds.json` 文件，填入公网代理服务器参数：

```json
{
  "outbounds": [
    {
      "tag": "nodeA",
      "...": "..."
    },
    {
      "tag": "nodeB",
      "...": "..."
    },
    {
      "tag": "nodeC",
      "...": "..."
    }
  ]
}
```

接着配置路由部分，让暴露的三个 socks5 接口对接到三台服务器上，并分别配置 IPv4 与 IPv6 的负载均衡；路由核心在这里接管所有流量，IPv4 流量应将随机转发到三台服务器，而 IPv6 流量只送往 `nodeC` 服务器；创建 `config/routing.json` 文件，写入以下配置：

```json
{
  "routing": {
    "domainStrategy": "AsIs",
    "rules": [
      {
        "type": "field",
        "inboundTag": ["nodeA"],
        "outboundTag": "nodeA"
      },
      {
        "type": "field",
        "inboundTag": ["nodeB"],
        "outboundTag": "nodeB"
      },
      {
        "type": "field",
        "inboundTag": ["nodeC"],
        "outboundTag": "nodeC"
      },
      {
        "type": "field",
        "ip": ["0.0.0.0/0"],
        "balancerTag": "ipv4"
      },
      {
        "type": "field",
        "ip": ["::/0"],
        "balancerTag": "ipv6"
      }
    ],
    "balancers": [
      {
        "tag": "ipv4",
        "selector": [ "nodeA", "nodeB", "nodeC" ]
      },
      {
        "tag": "ipv6",
        "selector": [ "nodeC" ]
      }
    ]
  }
}
```

重启 XProxy 容器使配置生效：

```bash
docker restart scutweb
```

最后，验证代理服务是否正常工作，若出现问题可以查看 `/etc/scutweb/log` 文件夹下的日志，定位错误原因。

## 代理 ICMP 流量

> 这一步仅用于修复 ICMP 代理，无此需求可以忽略。

由于 socks5 代理服务不支持 ICMP 协议，当前搭建的网络只有 TCP 与 UDP 发往外网，即使在上文我们注入了一段命令用于劫持 PING 流量，但是返回的仅仅是虚假结果，并没有实际意义；所以如果对这个缺陷不满，您可以考虑使用以下方法修复这个问题。

为了代理 ICMP 流量，我们必须选择网络层的 VPN 工具，从简单轻量可用方面考虑，`WireGuard` 比较适合当前应用场景：TCP 与 UDP 流量走 VLESS + XTLS 代理，ICMP 流量进入 WireGuard ，而 WireGuard 本身使用 UDP 协议传输，这些数据包通过 Xray 隧道再次代理送到远端服务器，解开后将 ICMP 流量送至公网；这种方式虽然略显繁杂，但实际场景中 ICMP 流量很少且数据包不大，并不存在性能问题。

具体实现上，我们需要在容器中安装 WireGuard 工具包，然后在 XProxy 中配置启动注入脚本，开启 WireGuard 对 ICMP 流量的代理。

### 1. 拉取 WireGuard 安装包

XProxy 容器默认不自带 WireGuard 功能，需要额外安装 `wireguard-tools` 包，您可以在原有镜像上添加一层，或是使用以下方式安装离线包。

> 以下代码用于生成 Alpine 的 WireGuard 安装脚本，您也可以选择手动拉取 apk 安装包

```python
#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import re
import sys

alpine = '3.16'

if len(sys.argv) != 2:
    print('Invalid argument')
    sys.exit(1)

workDir = os.path.join(os.path.dirname(os.path.realpath(__file__)), sys.argv[1])

output = os.popen(' '.join([
    'docker', 'run', '--rm', '-it', '-v', '%s:/pkg' % workDir, '-w', '/pkg',
    'alpine:%s' % alpine, 'sh', '-c', '\'apk update && apk fetch -R %s\'' % sys.argv[1]
])).read()
print("%(line)s\n%(msg)s%(line)s" % {'line': '=' * 88, 'msg': output})

with open(os.path.join(workDir, 'setup'), 'w') as script:
    script.write("#!/usr/bin/env sh\ncd \"$(dirname \"$0\")\"\napk add " + ' '.join([
        s + '.apk' for s in re.findall(r'Downloading (\S+)', output)
    ]) + " --no-network --quiet\n")
os.system('chmod +x %s' % os.path.join(workDir, 'setup'))
```

```bash
# fetch.py 为上述脚本
$ cd /etc/scutweb
$ mkdir -p ./toolset && cd ./toolset
$ python3 fetch.py wireguard-tools  # 拉取wireguard-tools依赖
···
···
```

拉取成功后将生成 `wireguard-tools` 文件夹，包含多个依赖的 `.apk` 安装包与 `setup` 安装脚本。

### 2. 写入 WireGuard 配置文件

一个典型的客户端配置文件如下：

```ini
[Interface]
PrivateKey = 客户端私钥

[Peer]
PublicKey = 服务端公钥
Endpoint = 服务器IP:端口
AllowedIPs = 0.0.0.0/0
```

将其保存至 `/etc/scutweb/config/wg.conf`

### 3. 容器注入 WireGuard 服务

WireGuard 在这里使用 `192.168.1.0/24` 的 VPN 网段，客户端 IP 地址为 `192.168.1.2`，注意服务端应允许 `192.168.2.2/24` 网段，否则必须在容器中多做一层 NAT 才能代理。

此外，XProxy 默认没有加入网关自身的代理，只在 `PREROUTING` 链上劫持流量，因此这里需要修改 `OUTPUT` 链，让 WireGuard 的流量被 XProxy 代理；将发往 WireGurad 服务器的流量打上标志 `0x1`，该数据包就会被重新路由到 `PREROUTING` 链上（netfilter 特性），从而进行透明代理。

```yaml
custom:
  pre:
    - /xproxy/toolset/wireguard-tools/setup  # 安装离线包
    - ip link add wg0 type wireguard
    - wg setconf wg0 /xproxy/config/wg.conf  # 加载配置文件
    - ip addr add 192.168.1.2/24 dev wg0  # 添加本机WireGuard地址
    - ip link set mtu 1420 up dev wg0  # 启动VPN服务
    - ip rule add fwmark 51820 table 51820
    - ip route add 0.0.0.0/0 dev wg0 table 51820  # WireGuard路由表
    - iptables -t mangle -N WGPROXY
    - iptables -t mangle -A WGPROXY -d 127.0.0.0/8 -j RETURN
    - iptables -t mangle -A WGPROXY -d 192.168.2.0/24 -j RETURN
    - iptables -t mangle -A WGPROXY -d 169.254.0.0/16 -j RETURN
    - iptables -t mangle -A WGPROXY -d 224.0.0.0/3 -j RETURN
    - iptables -t mangle -A WGPROXY -p icmp -j MARK --set-mark 51820  # ICMP流量送至WireGuard路由表
    - iptables -t mangle -A PREROUTING -j WGPROXY
    - iptables -t mangle -A OUTPUT -p udp -d 服务器IP –-dport 服务器端口 -j MARK --set-mark 1  # 重定向到PREROUTING
```

配置完成后，重启 XProxy 容器生效，在内网设备上执行 PING 命令，如果返回正常延迟则配置成功。
