package query

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

// parseHeader parses the DNS header from a byte slice
func ParseHeader(reader *bytes.Reader) DNSHeader {
	var header DNSHeader
	if err := binary.Read(reader, binary.BigEndian, &header); err != nil {
		log.Fatalf("binary.Read failed: %v", err)
	}
	return header
}

func DecodeName(reader *bytes.Reader) ([]byte, error) {
	var parts [][]byte // Use a slice of byte slices to collect parts of the name
	for {
		if reader.Len() == 0 {
			return nil, fmt.Errorf("premature end of data")
		}

		lengthByte, err := reader.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("reading length byte: %w", err)
		}

		if lengthByte == 0 {
			break // End of the name part, typically a zero byte indicates termination
		}

		if lengthByte&0xC0 == 0xC0 { // Check if it's a pointer (compression)
			nextByte, err := reader.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("reading second byte of pointer: %w", err)
			}
			pointer := (int(lengthByte&0x3F) << 8) | int(nextByte)
			currentPos, _ := reader.Seek(0, io.SeekCurrent) // Store current position

			_, err = reader.Seek(int64(pointer), io.SeekStart) // Jump to the pointer location
			if err != nil {
				return nil, fmt.Errorf("seeking to pointer location: %w", err)
			}

			part, err := DecodeName(reader) // Recursively decode the name at the pointer location
			if err != nil {
				return nil, err
			}
			parts = append(parts, part)

			_, err = reader.Seek(currentPos, io.SeekStart) // Restore the reader's position
			if err != nil {
				return nil, fmt.Errorf("restoring position after pointer: %w", err)
			}
			break // After decoding a compressed name, it ends the name field
		} else {
			part := make([]byte, lengthByte)
			_, err = reader.Read(part)
			if err != nil {
				return nil, fmt.Errorf("reading label: %w", err)
			}
			parts = append(parts, part)
		}
	}

	return bytes.Join(parts, []byte(".")), nil // Join parts with a dot as separator in byte form
}

func DecodeNSName(buffer []byte) string {
	domainName := ""
	pos := 0

	for pos < len(buffer) {
		segmentLength := int(buffer[pos])
		pos++
		if segmentLength == 0 {
			break
		}
		if len(domainName) > 0 {
			domainName += "."
		}
		domainName += string(buffer[pos : pos+segmentLength])
		pos += segmentLength
	}
	return domainName
}

func ParseQuestion(reader *bytes.Reader) (*DNSQuestion, error) {

	// Decode the DNS name
	name, err := DecodeName(reader)
	if err != nil {
		return &DNSQuestion{}, fmt.Errorf("failed to decode name: %w", err)
	}

	// Read the next four bytes for type and class
	typeClass := make([]byte, 4)
	if _, err := reader.Read(typeClass); err != nil {
		return &DNSQuestion{}, fmt.Errorf("failed to read type and class: %w", err)
	}

	// Unpack type and class
	var type_, class_ uint16
	if err := binary.Read(bytes.NewReader(typeClass[0:2]), binary.BigEndian, &type_); err != nil {
		return &DNSQuestion{}, fmt.Errorf("failed to unpack type: %w", err)
	}
	if err := binary.Read(bytes.NewReader(typeClass[2:]), binary.BigEndian, &class_); err != nil {
		return &DNSQuestion{}, fmt.Errorf("failed to unpack class: %w", err)
	}

	// print("Name: ", name)

	return &DNSQuestion{
		Name:  name,
		Type:  type_,
		Class: class_,
	}, nil
}

func ParseRecord(reader *bytes.Reader) (*DNSRecord, error) {
	name, err := DecodeName(reader) // Assume decodeName is defined correctly
	if err != nil {
		return &DNSRecord{}, fmt.Errorf("failed to decode name: %w", err)
	}
	var recordReader RecordReader

	if err := binary.Read(reader, binary.BigEndian, &recordReader); err != nil {
		return &DNSRecord{}, fmt.Errorf("failed to read type: %w", err)
	}
	length := make([]byte, recordReader.DataLen)
	if _, err := reader.Read(length); err != nil {
		return &DNSRecord{}, fmt.Errorf("failed to read data: %w", err)
	}

	return &DNSRecord{
		Name:  name,
		Type:  recordReader.Type,
		Class: recordReader.Class,
		TTL:   recordReader.TTL,
		Data:  length,
	}, nil
}
