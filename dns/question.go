package dns

import (
	"bytes"
	"encoding/binary"
	"log"
	"strings"
)

// Question type
// https://tools.ietf.org/html/rfc1035#section-4.1.2
type Question struct {
	QNAME string
	QTYPE uint16
	QCLASS uint16
}

// Encode method
func (q Question) Encode() []byte {
	var buffer bytes.Buffer

	domainParts := strings.Split(q.QNAME, ".")

	for _, part := range domainParts {
		if err := binary.Write(&buffer, binary.BigEndian, byte(len(part))); err != nil {
			log.Fatalf("Error binary.Write(..) for '%s': '%s'", part, err)
		}

		for _, char := range part {
			if err := binary.Write(&buffer, binary.BigEndian, uint8(char)); err != nil {
				log.Fatalf("Error binary.Write(..) for '%s'; '%c': '%s'", part, char, err)
			}
		}
	}

	binary.Write(&buffer, binary.BigEndian, uint8(0))
	binary.Write(&buffer, binary.BigEndian, q.QTYPE)
	binary.Write(&buffer, binary.BigEndian, q.QCLASS)

	return buffer.Bytes()
}
