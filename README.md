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
