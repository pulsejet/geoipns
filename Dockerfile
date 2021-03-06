# Build image
FROM golang:1.13-alpine
EXPOSE 5312/udp

WORKDIR /go/src/github.com/pulsejet/geoipns/

# Get dep
RUN apk add dep git bash zip unzip

# Copy files
COPY . .

# Build golang binary
RUN dep ensure && \
    go build

# Get GeoIP2 data
RUN ./getdata.sh && \
    zip -9 -r data.zip data

# Deploy image
FROM alpine:3.10
WORKDIR /geoipns
RUN apk add curl unzip

# Get binary and config
COPY --from=0 \
    /go/src/github.com/pulsejet/geoipns/geoipns \
    /go/src/github.com/pulsejet/geoipns/geoip/test/intranet.csv \
    /go/src/github.com/pulsejet/geoipns/data.zip \
    /go/src/github.com/pulsejet/geoipns/config.json \
    ./

# Fix configuration
RUN sed -i 's#geoip/test/intranet.csv#intranet.csv#g' config.json && \
    sed -i 's#127.0.0.1#0.0.0.0#g' config.json && \
    sed -i 's#"debug": true#"debug": false#g' config.json

# Environment
ENV SUFFIX=geoipns
ENV INTRANET_CSV_URL=https://raw.githubusercontent.com/pulsejet/geoipns/master/geoip/test/intranet.csv

# Execute
CMD curl -s -m 10 -o intranet.csv $INTRANET_CSV_URL && \
    unzip data.zip && rm data.zip && \
    sed -i "s#\"suffix\": \"geoipns\"#\"suffix\": \"$SUFFIX\"#g" config.json && \
    ./geoipns
