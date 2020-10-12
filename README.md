# GeoIPNS

Make DNS queries to get the geolocation and ASN name of an IP address. You can read the original blog post about why this project exists [here](https://medium.com/@varunpp2/locating-ip-addresses-with-dns-queries-af54228ea29c).

[![Build Status](https://travis-ci.org/pulsejet/geoipns.svg?branch=master)](https://travis-ci.org/pulsejet/geoipns)
[![codecov](https://codecov.io/gh/pulsejet/geoipns/branch/master/graph/badge.svg)](https://codecov.io/gh/pulsejet/geoipns)
[![Go Report Card](https://goreportcard.com/badge/github.com/pulsejet/geoipns)](https://goreportcard.com/report/github.com/pulsejet/geoipns)
[![GitHub license](https://img.shields.io/github/license/pulsejet/geoipns)](https://github.com/pulsejet/geoipns/blob/master/LICENSE)

## Usage

Make a `TXT` query for `IP.geoipns` to the server to get the location and ASN name of `IP` in an [RFC 1464](https://tools.ietf.org/html/rfc1464) compliant format. The model can be easily extended to store arbitrary data about IP subnets using the configuration file.

To build and run,
```shell
go get github.com/pulsejet/geoipns && cd $GOPATH/src/github.com/pulsejet/geoipns
dep ensure  # you may skip this step if deps are obtained automatically
go build
./getdata.sh  # write your custom script to get GeoIP data
./geoipns
```

To query,
```shell
$ dig +short @localhost -p5312 -tTXT 3.105.177.255.geoipns
"location=Sydney, NSW, AU"
"asn=Amazon.com, Inc."
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

## Benchmark
Testing GeoIPNS running on an LXC container limited to 4 vCPUs and 4GB RAM on a `24 x Intel(R) Xeon(R) CPU E5-2620 v2 @ 2.10GHz` machine with PHP7.3 on the same machine for 10000 randomized IPv4 requests gave the following results:
```text
Average 433μs per GeoIPNS request (native DNS client + PHP, different container)
Average 229μs per GeoIPNS request (phpdns, different container)
Average 150μs per GeoIPNS request (singe-threaded C++, same container)
Average 94μs per GeoIPNS request (goroutines, same container)
```

## Docker

> Since the GeoLite database cannot be downloaded directly any more, you will have to write your own script replacing getdata.sh to download the GeoIP database as a CSV. The `getdata.sh` and `Dockerfile` files are kept in the repository for historical purposes only.

Two environment variables can be set:
* `SUFFIX`: the DNS suffix to listen on
* `INTRANET_CSV_URL`: a URL pointing to a CSV file with location entries other than GeoIP2 lite
