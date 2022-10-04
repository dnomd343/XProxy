ARG ALPINE="alpine:3.16"
ARG GOLANG="golang:1.19-alpine3.16"

FROM ${ALPINE} AS upx
ENV UPX_VER="3.96"
RUN sed -i 's/v3.\d\d/v3.15/' /etc/apk/repositories && apk add bash build-base perl ucl-dev zlib-dev
RUN wget https://github.com/upx/upx/releases/download/v${UPX_VER}/upx-${UPX_VER}-src.tar.xz && tar xf upx-${UPX_VER}-src.tar.xz
WORKDIR ./upx-${UPX_VER}-src/
RUN make -C ./src/ && mkdir -p /upx/bin/ && mv ./src/upx.out /upx/bin/upx && \
    mkdir -p /upx/lib/ && cd /usr/lib/ && cp -d ./libgcc_s.so* ./libstdc++.so* ./libucl.so* /upx/lib/

FROM ${GOLANG} AS xray
ENV XRAY="1.6.0"
RUN wget https://github.com/XTLS/Xray-core/archive/refs/tags/v${XRAY}.tar.gz && tar xf v${XRAY}.tar.gz
WORKDIR ./Xray-core-${XRAY}/main/
RUN go get -d
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w" && mv main /tmp/xray

FROM ${GOLANG} AS v2ray
ENV V2FLY="5.1.0"
RUN wget https://github.com/v2fly/v2ray-core/archive/refs/tags/v${V2FLY}.tar.gz && tar xf v${V2FLY}.tar.gz
WORKDIR ./v2ray-core-${V2FLY}/main/
RUN go get -d
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w" && mv main /tmp/v2ray

FROM ${GOLANG} AS sagray
#ENV SAGER_VER="5.0.17"
RUN wget https://github.com/SagerNet/v2ray-core/archive/refs/heads/main.zip && unzip main.zip
WORKDIR ./v2ray-core-main/main/
#RUN wget https://github.com/SagerNet/v2ray-core/archive/refs/tags/v${SAGER}.tar.gz && tar xf v${SAGER}.tar.gz
#WORKDIR ./v2ray-core-${SAGER}/main/
RUN go get -d
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w" && mv main /tmp/sagray

FROM ${GOLANG} AS xproxy
COPY ./ /XProxy/
WORKDIR /XProxy/cmd/
RUN go get -d
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w" && mv cmd /tmp/xproxy

FROM ${ALPINE} AS build
WORKDIR /tmp/
RUN wget "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat" && \
    wget "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat"
RUN apk add xz && mkdir -p /asset/ && tar cJf /asset/assets.tar.xz ./*.dat
COPY --from=xproxy /tmp/xproxy /asset/usr/bin/
COPY --from=sagray /tmp/sagray /asset/usr/bin/
COPY --from=v2ray /tmp/v2ray /asset/usr/bin/
COPY --from=xray /tmp/xray /asset/usr/bin/
COPY --from=upx /upx/ /usr/
RUN ls /asset/usr/bin/* | xargs -P0 -n1 upx -9

FROM ${ALPINE}
RUN apk add --no-cache dhcp iptables ip6tables radvd && \
    mkdir -p /run/radvd/ && rm -f /etc/dhcp/dhcpd.conf.example && \
    touch /var/lib/dhcp/dhcpd.leases && touch /var/lib/dhcp/dhcpd6.leases
COPY --from=build /asset/ /
WORKDIR /xproxy/
ENTRYPOINT ["xproxy"]
