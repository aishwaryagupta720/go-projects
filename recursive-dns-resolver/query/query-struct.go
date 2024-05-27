package query

type DNSHeader struct {
	ID             uint16
	Flags          uint16
	NumQuestions   uint16
	NumAnswers     uint16
	NumAuthorities uint16
	NumAdditionals uint16
}

type DNSQuestion struct {
	Name  []byte
	Type  uint16
	Class uint16
}

// DNSRecord represents a DNS resource record
type DNSRecord struct {
	Name  []byte
	Type  uint16
	Class uint16
	TTL   uint32 //to accomodate bigger value
	Data  []byte
}

type RecordReader struct {
	Type    uint16
	Class   uint16
	TTL     uint32
	DataLen uint16
}

type DNSPacket struct {
	Header      DNSHeader
	Questions   []DNSQuestion
	Answers     []DNSRecord
	Authorities []DNSRecord
	Additionals []DNSRecord
}

const (
	TYPE_A     uint16 = 1
	TYPE_NS    uint16 = 2
	TYPE_CNAME uint16 = 5
	TYPE_TXT   uint16 = 16
	TYPE_AAAA  uint16 = 28
)
