package resolver

import (
	"fmt"
	"log"
	"net"
	"recursive-dns-resolver/cache"
	"recursive-dns-resolver/query"
	"recursive-dns-resolver/socket"
)

func SendQuery(domainName string, recordType uint16, root string) (*query.DNSPacket, error) {

	dnsquery := query.BuildQuery(domainName, recordType)
	// fmt.Printf("DNS Query: %x\n", dnsquery)

	if recordType == uint16(query.TYPE_AAAA) {

		root = "[" + root + "]"

	}

	// Send the DNS query
	conn, err := socket.OpenUDPConnection(root)
	if err != nil {
		log.Fatalf("Error opening socket: %v", err)
		return nil, err
	}
	defer socket.CloseUDPConnection(conn)

	_, err = conn.Write(dnsquery)
	if err != nil {
		log.Fatalf("Error sending query: %v", err)
		return nil, err
	}
	// fmt.Printf("DNS Request sent to %s for %s\n", root, domainName)
	// Buffer to receive the response
	buffer := make([]byte, 2048) // Sufficient size to handle typical DNS responses
	_, err = conn.Read(buffer)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}
	// fmt.Printf("Received DNS response (%d bytes): %x\n", n, buffer[:n])

	dnspacket, err := query.ParseDNSResponse(buffer)
	if err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}
	return dnspacket, err

}

func ResolveQuery(domainName string, recordType uint16) (string, error) {

	// rootservers := []string{"198.41.0.4", "199.9.14.201", "192.33.4.12", "199.7.91.13", "192.203.230.10", "192.5.5.241", "192.112.36.4", "198.97.190.53", "192.36.148.17", "192.58.128.30", "193.0.14.129", "199.7.83.42", "202.12.27.33"}
	// index := rand.Intn(13)
	root := "192.5.5.241"
	var original_record uint16
	var nsName string
	var nsIP string
	if recordType == uint16(query.TYPE_NS) || recordType == uint16(query.TYPE_CNAME) {
		original_record = recordType
		recordType = uint16(query.TYPE_A)
	}
	if original_record == uint16(query.TYPE_A) {
		cachedIP := cache.GetFromCache(domainName)
		if cachedIP != nil {
			return cachedIP.String(), nil
		}
	}

	for {
		response, err := SendQuery(domainName, recordType, root)
		if err != nil {
			return "", err
		}
		if ip, ttl := query.GetAnswerIP(*response); ip != "" {
			if original_record == uint16(query.TYPE_NS) {
				return nsName + " " + nsIP, nil
			}
			// if original_record == uint16(query.TYPE_CNAME) {
			// 	return query.GetAnswerCNAME(*response, uint16(query.TYPE_CNAME)), nil
			// }
			if original_record == uint16(query.TYPE_A) {
				cache.InsertInCache(domainName, net.ParseIP(ip), ttl)
			}
			return ip, nil
		} else if nsIP, nsName, _ = query.GetAdditionalsIP(*response, recordType); nsIP != "" {
			root = nsIP
		} else if nsDomain := query.GetNameServers(*response); nsDomain != "" {
			ip, err := ResolveQuery(nsDomain, recordType)
			if err != nil {
				return "", err
			}
			root = ip
		} else {
			return "", fmt.Errorf("something went wrong")
		}
	}
}
