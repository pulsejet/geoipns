# GeoIP NS

Make DNS queries to get the geolocation and ASN name of an IP address

[![Build Status](https://travis-ci.org/pulsejet/geoipns.svg?branch=master)](https://travis-ci.org/pulsejet/geoipns)
[![codecov](https://codecov.io/gh/pulsejet/geoipns/branch/master/graph/badge.svg)](https://codecov.io/gh/pulsejet/geoipns)
[![GitHub license](https://img.shields.io/github/license/pulsejet/geoipns)](https://github.com/pulsejet/geoipns/blob/master/LICENSE)

## Usage

Make a `TXT` query for `IP.geoipns` to the server to get the location and ASN name of `IP` in an [RFC 1464](https://tools.ietf.org/html/rfc1464) compliant format. The model can be easily extended to store arbitrary data about IP subnets using the configuration file.

To build and run
```
go get github.com/pulsejet/geoipns && cd $GOPATH/src/github.com/pulsejet/geoipns
dep ensure
go build
./getdata.sh
./geoipns
```

To query
```
dig @localhost -p5312 -tTXT 3.105.177.255.geoipns
```

will return

```
; <<>> DiG 9.11.3-1ubuntu1.7-Ubuntu <<>> @localhost -p5312 -tTXT 3.105.177.255.geoipns
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 40194
;; flags: qr rd; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;3.105.177.255.geoipns.         IN      TXT

;; ANSWER SECTION:
3.105.177.255.geoipns.  3600    IN      TXT     "location=Sydney, NSW, AU"
3.105.177.255.geoipns.  3600    IN      TXT     "asn=Amazon.com, Inc."

;; Query time: 0 msec
;; SERVER: 127.0.0.1#5312(127.0.0.1)
;; WHEN: Mon Oct 21 19:13:20 IST 2019
;; MSG SIZE  rcvd: 151
```

## Testing
Use `go test -v ./...` to run automated tests.
