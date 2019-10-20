# geolocation-dns

Make TXT queries to get the geolocation and ASN name of an IP address

To build
```
go get github.com/miekg/dns
go build
```

To query
```
dig @localhost -p5312 -tTXT 103.21.125.84.location
```
