package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

/*
	Example header:

	AA AA - ID
	01 00 - Query parameters (QR | Opcode | AA | TC | RD | RA | Z | ResponseCode)
	00 01 - Number of questions
	00 00 - Number of answers
	00 00 - Number of authority records
	00 00 - Number of additional records
*/

type DNSQuery struct {
	ID     uint16 // An arbitary 16 bit request identifier (same id is used in the response)
	QR     bool   // A 1 bit flat specifying whether this message is a query (0) or a response (1)
	Opcode uint8  // A 4 bit fields that specifies the query type; 0 (standard), 1 (inverse), 2 (status), 4 (notify), 5 (update)

	AA           bool  // Authoriative answer
	TC           bool  // 1 bit flag specifying if the message has been truncated
	RD           bool  // 1 bit flag to specify if recursion is desired (if the DNS server we secnd out request to doesn't know the answer to our query, it can recursively ask other DNS servers)
	RA           bool  // Recursive available
	Z            uint8 // Reserved for future use
	ResponseCode uint8

	QDCount uint16 // Number of entries in the question section
	ANCount uint16 // Number of answers
	NSCount uint16 // Number of authorities
	ARCount uint16 // Number of additional records

	Questions []DNSQuestion
}

func (q DNSQuery) encode() []byte {

	q.QDCount = uint16(len(q.Questions))

	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, q.ID)

	b2i := func(b bool) int {
		if b {
			return 1
		}

		return 0
	}

	queryParams1 := byte(b2i(q.QR)<<7 | int(q.Opcode)<<3 | b2i(q.AA)<<1 | b2i(q.RD))
	queryParams2 := byte(b2i(q.RA)<<7 | int(q.Z)<<4)

	binary.Write(&buffer, binary.BigEndian, queryParams1)
	binary.Write(&buffer, binary.BigEndian, queryParams2)
	binary.Write(&buffer, binary.BigEndian, q.QDCount)
	binary.Write(&buffer, binary.BigEndian, q.ANCount)
	binary.Write(&buffer, binary.BigEndian, q.NSCount)
	binary.Write(&buffer, binary.BigEndian, q.ARCount)

	for _, question := range q.Questions {
		buffer.Write(question.encode())
	}

	return buffer.Bytes()
}

/*
	Example Question:

	07 65 - 'example' has length 7, e
	78 61 - x, a
	6D 70 - m, p
	6C 65 - l, e
	03 63 - 'com' has length 3, c
	6F 6D - o, m
	00    - zero byte to end the QNAME
	00 01 - QTYPE
	00 01 - QCLASS

	76578616d706c6503636f6d0000010001
*/

type DNSQuestion struct {
	Domain string
	Type   uint16 // DNS Record type we are looking up; 1 (A record), 2 (authoritive name server)
	Class  uint16 // 1 (internet)
}

func (q DNSQuestion) encode() []byte {
	var buffer bytes.Buffer

	domainParts := strings.Split(q.Domain, ".")
	for _, part := range domainParts {
		if err := binary.Write(&buffer, binary.BigEndian, byte(len(part))); err != nil {
			log.Fatalf("Error binary.Write(..) for '%s': '%s'", part, err)
		}

		for _, c := range part {
			if err := binary.Write(&buffer, binary.BigEndian, uint8(c)); err != nil {
				log.Fatalf("Error binary.Write(..) for '%s'; '%c': '%s'", part, c, err)
			}
		}
	}

	binary.Write(&buffer, binary.BigEndian, uint8(0))
	binary.Write(&buffer, binary.BigEndian, q.Type)
	binary.Write(&buffer, binary.BigEndian, q.Class)

	return buffer.Bytes()

}

func printResponseCode(responseCode byte) {

	switch responseCode {
	case 0:
		fmt.Println("Domain exists!")
	case 1:
		fmt.Println("Format error")
	case 2:
		fmt.Println("Server failure")
	case 3:
		fmt.Println("Non-existent domain")
	case 9:
		fmt.Println("Server not authorative for zone")
	case 10:
		fmt.Println("Name not in zone")
	default:
		fmt.Printf("Unmapped response code for '%d'\n", responseCode)
	}
}

func main() {
	host := "8.8.8.8:53"
	domain := "example.com"

	q := DNSQuestion{
		Domain: domain,
		Type:   0x1, // A record
		Class:  0x1, // Internet
	}

	query := DNSQuery{
		ID:        0xAAAA,
		RD:        true,
		Questions: []DNSQuestion{q},
	}

	// Setup a UDP connection
	conn, err := net.Dial("udp", host)
	if err != nil {
		log.Fatal("failed to connect:", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		log.Fatal("failed to set deadline: ", err)
	}

	encodedQuery := query.encode()

	conn.Write(encodedQuery)

	encodedAnswer := make([]byte, len(encodedQuery))
	if _, err := bufio.NewReader(conn).Read(encodedAnswer); err != nil {
		log.Fatal(err)
	}

	responseCode := encodedAnswer[3] & 0xF

	fmt.Printf("DNS query results for '%s' against '%s'\n", domain, host)
	fmt.Printf(">> ")
	printResponseCode(responseCode)
	fmt.Printf(">> encodedQuery:  %#x\n", encodedQuery)
	fmt.Printf(">> encodedAnswer: %#x\n", encodedAnswer)
}
