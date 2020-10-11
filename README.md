# Manual DNS Lookup

## Description
Original part:
> This is an example of how to query a DNS server (in this example "8.8.8.8") for "A" records using handcrafted packets.  
This shows the raw bytes to send to a DNS server in order to receive an answer.

I forked this repository to practice Go and network requests at a lower level. I plan to implement decoding of the response from the DNS server and remove the input parameters from the code.

## Running

```
$ go run main.go 
DNS query results for 'example.com' against '8.8.8.8:53'
>> Domain exists!
>> encodedQuery:  0xaaaa01000001000000000000076578616d706c6503636f6d0000010001
>> encodedAnswer: 0xaaaa81800001000100000000076578616d706c6503636f6d0000010001
```

## Resources
- [Original project](https://github.com/vishen/go-dnsquery)
- [RFC 1035 - Domain names - implementation and specification](https://tools.ietf.org/html/rfc1035)
