ARG ALPINE="alpine:3.16"
ARG GOLANG="golang:1.19-alpine3.16"

FROM ${ALPINE} AS upx
RUN apk add build-base cmake git
RUN git clone https://github.com/dnomd343/upx.git
WORKDIR ./upx/
RUN git submodule update --init && rm -rf ./.git/
RUN make UPX_CMAKE_CONFIG_FLAGS=-DCMAKE_EXE_LINKER_FLAGS=-static && \
    mv ./build/release/upx /tmp/ && strip /tmp/upx

FROM ${GOLANG} AS xray
ENV XRAY="1.6.0"
RUN wget https://github.com/XTLS/Xray-core/archive/refs/tags/v${XRAY}.tar.gz && tar xf v${XRAY}.tar.gz
WORKDIR ./Xray-core-${XRAY}/main/
RUN go get -d
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w" && mv main /tmp/xray
COPY --from=upx /tmp/upx /usr/bin/
RUN upx -9 /tmp/xray

FROM ${GOLANG} AS v2ray
ENV V2FLY="5.1.0"
RUN wget https://github.com/v2fly/v2ray-core/archive/refs/tags/v${V2FLY}.tar.gz && tar xf v${V2FLY}.tar.gz
WORKDIR ./v2ray-core-${V2FLY}/main/
RUN go get -d
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w" && mv main /tmp/v2ray
COPY --from=upx /tmp/upx /usr/bin/
RUN upx -9 /tmp/v2ray

FROM ${GOLANG} AS sagray
#ENV SAGER_VER="5.0.17"
RUN wget https://github.com/SagerNet/v2ray-core/archive/refs/heads/main.zip && unzip main.zip
WORKDIR ./v2ray-core-main/main/
#RUN wget https://github.com/SagerNet/v2ray-core/archive/refs/tags/v${SAGER}.tar.gz && tar xf v${SAGER}.tar.gz
#WORKDIR ./v2ray-core-${SAGER}/main/
RUN go get -d
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w" && mv main /tmp/sagray
COPY --from=upx /tmp/upx /usr/bin/
RUN upx -9 /tmp/sagray

FROM ${GOLANG} AS xproxy
COPY ./ /XProxy/
WORKDIR /XProxy/cmd/
RUN go get -d
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w" && mv cmd /tmp/xproxy
COPY --from=upx /tmp/upx /usr/bin/
RUN upx -9 /tmp/xproxy

FROM ${ALPINE} AS build
RUN apk add xz
WORKDIR /release/
RUN wget "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat" && \
    wget "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat" && \
    tar cJf assets.tar.xz *.dat && rm *.dat
COPY --from=xproxy /tmp/xproxy /release/usr/bin/
COPY --from=sagray /tmp/sagray /release/usr/bin/
COPY --from=v2ray /tmp/v2ray /release/usr/bin/
COPY --from=xray /tmp/xray /release/usr/bin/

FROM ${ALPINE}
RUN apk add --no-cache dhcp iptables ip6tables radvd && \
    mkdir -p /run/radvd/ && rm -f /etc/dhcp/dhcpd.conf.example && \
    touch /var/lib/dhcp/dhcpd.leases && touch /var/lib/dhcp/dhcpd6.leases
COPY --from=build /release/ /
WORKDIR /xproxy/
ENTRYPOINT ["xproxy"]
