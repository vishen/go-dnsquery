package dns

import (
	"bytes"
	"encoding/binary"
)

// Query type
// https://tools.ietf.org/html/rfc1035#section-4.1.1
type Query struct {
	ID uint16
	QR bool
	OPCODE uint8
	AA bool
	TC bool
	RD bool
	RA bool
	Z uint8
	RCODE uint8
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
	Questions []Question
}

// Encode method
func (q Query) Encode() []byte {
	q.QDCOUNT = uint16(len(q.Questions))

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, q.ID)

	b2i := func (b bool) int {
		if b {
			return 1
		}
		return 0
	}

	queryParams1 := byte(b2i(q.QR)<<7 | int(q.OPCODE)<<3 | b2i(q.AA)<<1 | b2i(q.RD))
	queryParams2 := byte(b2i(q.RA)<<7 | int(q.Z)<<4)

	binary.Write(&buffer, binary.BigEndian, queryParams1)
	binary.Write(&buffer, binary.BigEndian, queryParams2)
	binary.Write(&buffer, binary.BigEndian, q.QDCOUNT)
	binary.Write(&buffer, binary.BigEndian, q.ANCOUNT)
	binary.Write(&buffer, binary.BigEndian, q.NSCOUNT)
	binary.Write(&buffer, binary.BigEndian, q.ARCOUNT)

	for _, question := range q.Questions {
		buffer.Write(question.Encode())
	}

	return buffer.Bytes()
}
