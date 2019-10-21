# GeoIP NS

Make DNS queries to get the geolocation and ASN name of an IP address

[![Build Status](https://travis-ci.org/pulsejet/geoipns.svg?branch=master)](https://travis-ci.org/pulsejet/geoipns)
[![codecov](https://codecov.io/gh/pulsejet/geoipns/branch/master/graph/badge.svg)](https://codecov.io/gh/pulsejet/geoipns)
[![GitHub license](https://img.shields.io/github/license/pulsejet/geoipns)](https://github.com/pulsejet/geoipns/blob/master/LICENSE)

## Usage

Make a `TXT` query for `ip.location` to the server to get the response in the format `city,_province,_country|my_asn_name|S`

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
dig @localhost -p5312 -tTXT 103.21.125.84.location
```

will return

```
; <<>> DiG 9.11.3-1ubuntu1.7-Ubuntu <<>> @localhost -p5312 -tTXT 103.21.125.84.location
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 3009
;; flags: qr rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;103.21.125.84.location.                IN      TXT

;; ANSWER SECTION:
103.21.125.84.location. 3600    IN      TXT     "Mumbai,_MH,_IN|Powai|S"

;; Query time: 0 msec
;; SERVER: 127.0.0.1#5312(127.0.0.1)
;; WHEN: Sun Oct 20 17:45:30 IST 2019
;; MSG SIZE  rcvd: 97
```

# Testing
Use `go test -v` in the `geoip` folder to run automated tests.
