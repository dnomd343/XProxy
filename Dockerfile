ARG ALPINE_IMG="alpine:3.16"
ARG GOLANG_IMG="golang:1.19-alpine3.16"

FROM ${ALPINE_IMG} AS upx
ENV UPX_VER="3.96"
RUN sed -i 's/v3.\d\d/v3.15/' /etc/apk/repositories && apk add bash build-base perl ucl-dev zlib-dev
RUN wget https://github.com/upx/upx/releases/download/v${UPX_VER}/upx-${UPX_VER}-src.tar.xz && tar xf upx-${UPX_VER}-src.tar.xz
WORKDIR ./upx-${UPX_VER}-src/
RUN make -C ./src/ && mkdir -p /upx/bin/ && mv ./src/upx.out /upx/bin/upx && \
    mkdir -p /upx/lib/ && cd /usr/lib/ && cp -d ./libgcc_s.so* ./libstdc++.so* ./libucl.so* /upx/lib/

FROM ${GOLANG_IMG} AS xray
ENV XRAY_VER="1.6.0"
RUN wget https://github.com/XTLS/Xray-core/archive/refs/tags/v${XRAY_VER}.tar.gz && tar xf v${XRAY_VER}.tar.gz
WORKDIR ./Xray-core-${XRAY_VER}/main/
RUN go get -d
RUN env CGO_ENABLED=0 go build -v -o xray -trimpath -ldflags "-s -w" && mv xray /tmp/

FROM ${GOLANG_IMG} AS v2ray
ENV V2FLY_VER="5.1.0"
RUN wget https://github.com/v2fly/v2ray-core/archive/refs/tags/v${V2FLY_VER}.tar.gz && tar xf v${V2FLY_VER}.tar.gz
WORKDIR ./v2ray-core-${V2FLY_VER}/main/
RUN go get -d
RUN env CGO_ENABLED=0 go build -v -o v2ray -trimpath -ldflags "-s -w" && mv v2ray /tmp/

FROM ${GOLANG_IMG} AS sagray
#ENV SAGER_VER="5.0.16"
#RUN wget https://github.com/SagerNet/v2ray-core/archive/refs/tags/v${SAGER_VER}.tar.gz && tar xf v${SAGER_VER}.tar.gz
#WORKDIR ./v2ray-core-${SAGER_VER}/main/
RUN apk add git && git clone https://github.com/SagerNet/v2ray-core.git
WORKDIR ./v2ray-core/main/
RUN go get -d
RUN env CGO_ENABLED=0 go build -v -o sagray -trimpath -ldflags "-s -w" && mv sagray /tmp/

FROM ${GOLANG_IMG} AS xproxy
COPY . /XProxy
WORKDIR /XProxy/cmd/
RUN go get -d
RUN env CGO_ENABLED=0 go build -v -o xproxy -trimpath -ldflags "-s -w" && mv xproxy /tmp/

FROM ${ALPINE_IMG} AS asset
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

FROM ${ALPINE_IMG}
RUN apk add --no-cache dhcp iptables ip6tables radvd && \
    mkdir -p /run/radvd/ && rm -f /etc/dhcp/dhcpd.conf.example && \
    touch /var/lib/dhcp/dhcpd.leases && touch /var/lib/dhcp/dhcpd6.leases
COPY --from=asset /asset/ /
WORKDIR /xproxy
ENTRYPOINT ["xproxy"]
