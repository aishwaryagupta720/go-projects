package query

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

func ParseDNSResponse(buffer []byte) (*DNSPacket, error) {
	reader := bytes.NewReader(buffer)

	header := ParseHeader(reader)
	// fmt.Println("Parsed Header: ", header)

	questions := make([]DNSQuestion, header.NumQuestions)
	for i := 0; i < int(header.NumQuestions); i++ {
		question, err := ParseQuestion(reader)
		if err != nil {
			log.Printf("Error parsing question: %v", err)
			continue
		}
		questions[i] = *question
		// fmt.Printf("Question %d: %+v\n", i+1, question)
	}
	answers := make([]DNSRecord, header.NumAnswers)
	for i := 0; i < int(header.NumAnswers); i++ {
		answer, err := ParseRecord(reader)
		if err != nil {
			log.Printf("Error parsing answer: %v", err)
			continue
		}
		answers[i] = *answer
		// fmt.Printf("Answer %d: %+v\n", i+1, answer)
	}
	authorities := make([]DNSRecord, header.NumAuthorities)
	for i := 0; i < int(header.NumAuthorities); i++ {
		authority, err := ParseRecord(reader)
		if err != nil {
			log.Printf("Error parsing authority record: %v", err)
			continue
		}
		authorities[i] = *authority
		// fmt.Printf("Authority %d: %+v\n", i+1, authority)
	}
	additionals := make([]DNSRecord, header.NumAdditionals)
	for i := 0; i < int(header.NumAdditionals); i++ {
		additional, err := ParseRecord(reader)
		if err != nil {
			log.Printf("Error parsing additional record: %v", err)
			continue
		}
		additionals[i] = *additional
		// fmt.Printf("Additional %d: %+v\n", i+1, additional)
	}
	return &DNSPacket{
		Header:      header,
		Questions:   questions,
		Answers:     answers,
		Authorities: authorities,
		Additionals: additionals,
	}, nil
}

// func GetHeaderTTL(records DNSPacket) {
// 	header := records.Header
// 	print(header)

// }
func GetAnswerIP(records DNSPacket) (string, uint32) {
	var ips []string
	var ttl uint32
	for _, record := range records.Answers {
		ip := net.IP(record.Data).String()
		ttl = record.TTL
		if ip != "<nil>" {
			ips = append(ips, ip)
		}
	}
	if len(ips) > 0 {
		return ips[0], ttl
	}
	return "", 0
}

func GetAdditionalsIP(records DNSPacket, recordType uint16) (string, string, error) {
	for _, record := range records.Additionals {
		if record.Type == recordType {
			nsname := string(record.Name)
			ip := net.IP(record.Data)
			if ip == nil || ip.String() == "<nil>" {
				return "", "", fmt.Errorf("invalid IP address data")
			}
			return ip.String(), nsname, nil
		}
	}
	return "", "", fmt.Errorf("no record found of type %d", recordType)
}

func GetNameServers(packet DNSPacket) string {
	var nsDomains []string
	for _, record := range packet.Authorities {
		if record.Type == TYPE_NS { // Compare record type
			nsName := DecodeNSName(record.Data)
			if nsName != "" {
				nsDomains = append(nsDomains, nsName)
			}
		}
	}

	// Return a single string joined by commas if multiple NS records exist
	if len(nsDomains) > 0 {
		return nsDomains[0]
	}
	return ""
}
