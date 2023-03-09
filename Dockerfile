ARG ALPINE="alpine:3.17"
ARG GOLANG="golang:1.20-alpine3.17"

FROM ${ALPINE} AS upx
RUN apk add build-base cmake
ENV UPX="4.0.2"
RUN wget https://github.com/upx/upx/releases/download/v${UPX}/upx-${UPX}-src.tar.xz && tar xf upx-${UPX}-src.tar.xz
WORKDIR ./upx-${UPX}-src/
RUN make UPX_CMAKE_CONFIG_FLAGS=-DCMAKE_EXE_LINKER_FLAGS=-static
WORKDIR ./build/release/
RUN strip upx && mv upx /tmp/

FROM ${GOLANG} AS xray
#ENV XRAY="1.7.5"
#RUN wget https://github.com/XTLS/Xray-core/archive/refs/tags/v${XRAY}.tar.gz && tar xf v${XRAY}.tar.gz
#WORKDIR ./Xray-core-${XRAY}/main/

# TODO: use xray dev version just for now
RUN apk add git
RUN git clone https://github.com/XTLS/Xray-core.git
WORKDIR ./Xray-core/main/
RUN git checkout 4c8ee0af50bbabd29e6766f0d9509add6fc0b2e7

RUN go get
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w" && mv main /tmp/xray
COPY --from=upx /tmp/upx /usr/bin/
RUN upx -9 /tmp/xray

FROM ${GOLANG} AS xproxy
RUN apk add git
COPY ./ /XProxy/
WORKDIR /XProxy/cmd/
RUN go get
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-X main.version=$(git describe --tag) -s -w" && mv cmd /tmp/xproxy
COPY --from=upx /tmp/upx /usr/bin/
RUN upx -9 /tmp/xproxy

FROM ${ALPINE} AS build
RUN apk add xz
WORKDIR /release/
RUN wget "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat" && \
    wget "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat" && \
    tar cJf assets.tar.xz *.dat && rm *.dat
COPY --from=xproxy /tmp/xproxy /release/usr/bin/
COPY --from=xray /tmp/xray /release/usr/bin/

FROM ${ALPINE}
RUN apk add --no-cache dhcp iptables ip6tables radvd && \
    cd /var/lib/dhcp/ && touch dhcpd.leases dhcpd6.leases && \
    rm -f /etc/dhcp/dhcpd.conf.example && mkdir -p /run/radvd/
COPY --from=build /release/ /
WORKDIR /xproxy/
ENTRYPOINT ["xproxy"]
