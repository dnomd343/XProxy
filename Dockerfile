ARG ALPINE="alpine:3.19"
ARG GOLANG="golang:1.21-alpine3.19"

FROM ${GOLANG} AS xray
ENV XRAY="1.8.6"
RUN wget https://github.com/XTLS/Xray-core/archive/v${XRAY}.tar.gz -O- | tar xz
WORKDIR ./Xray-core-${XRAY}/main/
RUN go get
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w"
RUN mv main /xray

FROM ${GOLANG} AS xproxy
RUN apk add git
COPY ./ /XProxy/
WORKDIR /XProxy/cmd/
RUN go get
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-X main.version=$(git describe --tag) -s -w"
RUN mv cmd /xproxy

FROM ${ALPINE} AS assets
RUN apk add xz
RUN wget "https://cdn.dnomd343.top/v2ray-rules-dat/geoip.dat"
RUN wget "https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat"
RUN tar cJf /assets.tar.xz geoip.dat geosite.dat

FROM ${ALPINE} AS release
RUN apk add upx
WORKDIR /release/run/radvd/
WORKDIR /release/var/lib/dhcp/
RUN touch dhcpd.leases dhcpd6.leases
COPY --from=xproxy /xproxy /release/usr/bin/
COPY --from=assets /assets.tar.xz /release/
COPY --from=xray /xray /release/usr/bin/
WORKDIR /release/usr/bin/
RUN ls | xargs -n1 -P0 upx -9

FROM ${ALPINE}
RUN apk add --no-cache dhcp radvd iptables ip6tables
COPY --from=release /release/ /
WORKDIR /xproxy/
ENTRYPOINT ["xproxy"]
