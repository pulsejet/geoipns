# Build image
FROM golang:1.13-alpine
EXPOSE 5312/udp

WORKDIR /go/src/github.com/pulsejet/geoipns/

# Get dep
RUN apk add dep git bash

# Copy files
COPY . .

# Build golang binary
RUN dep ensure && \
    go build

# Get GeoIP2 data
RUN ./getdata.sh

# Deploy image
FROM alpine:3.10
WORKDIR /geoipns
RUN apk add curl

# Get binary and config
COPY --from=0 \
    /go/src/github.com/pulsejet/geoipns/geoipns \
    /go/src/github.com/pulsejet/geoipns/geoip/test/intranet.csv \
    /go/src/github.com/pulsejet/geoipns/config.json \
    ./

# Get data
COPY --from=0 \
    /go/src/github.com/pulsejet/geoipns/data \
    ./data

# Fix configuration
RUN sed -i 's#geoip/test/intranet.csv#intranet.csv#g' config.json && \
    sed -i 's#127.0.0.1#0.0.0.0#g' config.json && \
    sed -i 's#"debug": true#"debug": false#g' config.json

# Environment
ENV SUFFIX=geoipns
ENV INTRANET_CSV_URL=https://raw.githubusercontent.com/pulsejet/geoipns/master/geoip/test/intranet.csv

# Execute
CMD curl -s -m 10 -o intranet.csv $INTRANET_CSV_URL && \
    sed -i "s#\"suffix\": \"geoipns\"#\"suffix\": \"$SUFFIX\"#g" config.json && \
    ./geoipns
