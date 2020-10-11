// Original project: https://github.com/vishen/go-dnsquery
package main

import (
	"bufio"
	"./dns"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	host := "8.8.8.8:53"
	domain := "example.com"

	question := dns.Question {
		QNAME: domain,
		QTYPE: 0x1,
		QCLASS: 0x1,
	}

	query := dns.Query {
		ID: 0xAAAA,
		RD: true,
		Questions: []dns.Question{question},
	}

	conn, err := net.Dial("udp", host)

	if err != nil {
		log.Fatal("Failed to connect: ", err)
	}

	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		log.Fatal("Failed to set deadline: ", err)
	}

	encodedQuery := query.Encode()
	conn.Write(encodedQuery)

	encodedAnswer := make([]byte, len(encodedQuery))
	if _, err := bufio.NewReader(conn).Read(encodedAnswer); err != nil {
		log.Fatal(err)
	}

	responseCode := encodedAnswer[3] & 0xF

	fmt.Printf("DNS query results for '%s' against '%s'\n>> ", domain, host)
	printResponseCode(responseCode)
	fmt.Printf(">> encodedQuery:  %#x\n", encodedQuery)
	fmt.Printf(">> encodedAnswer: %#x\n", encodedAnswer)
}

func printResponseCode(responseCode byte) {
	switch responseCode {
	case 0: fmt.Println("Domain exists!")
	case 1: fmt.Println("Format error")
	case 2: fmt.Println("Server failure")
	case 3: fmt.Println("Non-existent domain")
	case 9: fmt.Println("Server not authorative for zone")
	case 10: fmt.Println("Name not in zone")
	default: fmt.Printf("Unmapped response code for '%d'\n", responseCode)
	}
}
