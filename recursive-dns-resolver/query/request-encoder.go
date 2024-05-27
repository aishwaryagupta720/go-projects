package query

import (
	"bytes"
	"encoding/binary"
	"math/rand"
)

// All fields in DNSHeader are of the same data type so it's efficient and safe to serialize the whole struct at once
func headerToBytes(header DNSHeader) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, header)
	return buf.Bytes()
}

// After writing the variable-length Name, the fixed-length fields Type and Class are written.
func questionToBytes(question DNSQuestion) []byte {
	buf := new(bytes.Buffer)
	buf.Write(question.Name)
	binary.Write(buf, binary.BigEndian, question.Type)
	binary.Write(buf, binary.BigEndian, question.Class)
	return buf.Bytes()
}

// encodeDNSName encodes a domain name into the DNS query format.
func encodeDNSName(domainName string) []byte {
	var buffer bytes.Buffer
	parts := bytes.Split([]byte(domainName), []byte("."))
	for _, part := range parts {
		if len(part) > 0 { // Ensure the part is non-empty
			buffer.WriteByte(byte(len(part))) // Length of the part
			buffer.Write(part)                // Actual part bytes
		}
	}
	buffer.WriteByte(0) // Null terminator for the domain name
	return buffer.Bytes()
}

// buildQuery constructs the DNS query byte sequence
func BuildQuery(domainName string, recordType uint16) []byte {
	name := encodeDNSName(domainName)
	id := uint16(rand.Intn(65536)) // Random ID from 0 to 65535
	recursionDesired := uint16(1 << 8)
	header := DNSHeader{
		ID:             id,
		Flags:          recursionDesired,
		NumQuestions:   1,
		NumAnswers:     0,
		NumAuthorities: 0,
		NumAdditionals: 0,
	}
	question := DNSQuestion{
		Name:  name,
		Type:  recordType,
		Class: 1, // CLASS_IN
	}
	query := append(headerToBytes(header), questionToBytes(question)...)
	return query
}
