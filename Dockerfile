ARG ALPINE="alpine:3.17"
ARG GOLANG="golang:1.20-alpine3.17"

FROM ${GOLANG} AS xray
ENV XRAY="1.8.0"
RUN wget https://github.com/XTLS/Xray-core/archive/refs/tags/v${XRAY}.tar.gz && tar xf v${XRAY}.tar.gz
WORKDIR ./Xray-core-${XRAY}/main/
RUN go get
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w" && mv main /tmp/xray

FROM ${GOLANG} AS xproxy
RUN apk add git
COPY ./ /XProxy/
WORKDIR /XProxy/cmd/
RUN go get
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-X main.version=$(git describe --tag) -s -w" && mv cmd /tmp/xproxy

FROM ${ALPINE} AS build
RUN apk add upx xz
WORKDIR /release/
RUN wget "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat" && \
    wget "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat" && \
    tar cJf assets.tar.xz *.dat && rm *.dat
COPY --from=xproxy /tmp/xproxy /release/usr/bin/
COPY --from=xray /tmp/xray /release/usr/bin/
WORKDIR /release/usr/bin/
RUN ls | xargs -n1 -P0 upx -9

FROM ${ALPINE}
RUN apk add --no-cache dhcp iptables ip6tables radvd && \
    cd /var/lib/dhcp/ && touch dhcpd.leases dhcpd6.leases && \
    rm -f /etc/dhcp/dhcpd.conf.example && mkdir -p /run/radvd/
COPY --from=build /release/ /
WORKDIR /xproxy/
ENTRYPOINT ["xproxy"]
