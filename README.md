# Manual DNS Lookup
This is an example of how to query a DNS server (in this example "8.8.8.8") for "A" records using handcrafted packets.

This shows the raw bytes to send to a DNS server in order to receive an answer.


## Running
```
$ go run main.go
DNS query results for 'example.com' against '8.8.8.8:53'
>> Domain exists!
>> encodedQuery:  0xaaaa01000001000000000000076578616d706c6503636f6d0000010001
>> encodedAnswer: 0xaaaa81800001000100000000076578616d706c6503636f6d0000010001
```

## Resources
```
- https://routley.io/tech/2017/12/28/hand-writing-dns-messages.html
```
