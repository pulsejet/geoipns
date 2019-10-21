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

## Data Structure
Each subnet is stored as an IP range with start and end 128-bit IP addresses. All data points are stored together as an array sorted by the IP address, with a high/low boolean. The data structure is looked up using binary search to find the location of the queried IP address. If a low IP precedes the discovered location, the IP is contained in the corresponding range; if a high IP precedes then the IP is contained in the parent range (if existent). Each data point stores information about the complementary point and the parent range.

## Example Data Structure
The following CSV (with proper configuration)

| cidr            | start_ip       | end_ip         | location     | 
|-----------------|----------------|----------------|--------------| 
| 10.105.0.0/16   |                |                | OuterIPRange | 
| 10.105.177.0/24 |                |                | Subnet1      | 
|                 | 10.105.177.120 | 10.105.177.200 | Subnet11     | 
|                 | 10.105.177.120 | 10.105.177.128 | Subnet12     | 
|                 | 10.105.200.0   | 10.105.201.255 | Subnet2      | 

would produce the following data structure (note that all IP addresses are internally represented as 128-bit bytes slices)

| Address | IP             | IsHigh | Data         | Complement | Parent | 
|---------|----------------|--------|--------------|------------|--------| 
| 0x001   | 10.105.0.0     | false  | OuterIPRange | 0x010      | 0x000  | 
| 0x002   | 10.105.177.0   | false  | Subnet1      | 0x007      | 0x001  | 
| 0x003   | 10.105.177.120 | false  | Subnet11     | 0x006      | 0x002  | 
| 0x004   | 10.105.177.120 | false  | Subnet12     | 0x005      | 0x003  | 
| 0x005   | 10.105.177.128 | true   | Subnet12     | 0x004      | 0x003  | 
| 0x006   | 10.105.177.200 | true   | Subnet11     | 0x003      | 0x002  | 
| 0x007   | 10.105.177.255 | true   | Subnet1      | 0x002      | 0x001  | 
| 0x008   | 10.105.200.0   | false  | Subnet2      | 0x009      | 0x001  | 
| 0x009   | 10.105.201.255 | true   | Subnet2      | 0x008      | 0x001  | 
| 0x010   | 10.105.255.255 | true   | OuterIPRange | 0x001      | 0x000  | 

## Testing
Use `go test -v ./...` to run automated tests.
