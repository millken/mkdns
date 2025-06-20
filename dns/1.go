package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// DNSHeader 表示DNS消息头
type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QdCount uint16
	AnCount uint16
	NsCount uint16
	ArCount uint16
}

// DNSQuestion 表示DNS问题部分
type DNSQuestion struct {
	Name  []byte
	Type  uint16
	Class uint16
}

// DNSRecord 表示DNS资源记录
type DNSRecord struct {
	Name  []byte
	Data  []byte
	TTL   uint32
	Type  uint16
	Class uint16
}

// 解析DNS消息
func parseDNSMessage(data []byte) {
	reader := bytes.NewReader(data)

	// 解析DNS消息头
	var header DNSHeader
	if err := binary.Read(reader, binary.BigEndian, &header); err != nil {
		fmt.Println("Error reading DNS header:", err)
		return
	}

	// 解析问题部分
	for i := 0; i < int(header.QdCount); i++ {
		question, err := parseDNSQuestion(reader)
		if err != nil {
			fmt.Println("Error parsing DNS question:", err)
			return
		}
		fmt.Printf("Question: %s, Type: %d, Class: %d\n", question.Name, question.Type, question.Class)
	}

	// 解析额外记录部分
	for i := 0; i < int(header.ArCount); i++ {
		record, err := parseDNSRecord(reader)
		if err != nil {
			fmt.Println("Error parsing DNS record:", err)
			return
		}
		if record.Type == 41 { // 41 是OPT记录的类型
			fmt.Println("Found OPT record")
			parseOPTData(record.Data)
		}
	}
}

// 解析DNS问题部分
func parseDNSQuestion(reader *bytes.Reader) (DNSQuestion, error) {
	var question DNSQuestion
	name, err := parseDNSName(reader)
	if err != nil {
		return question, err
	}
	question.Name = name
	if err := binary.Read(reader, binary.BigEndian, &question.Type); err != nil {
		return question, err
	}
	if err := binary.Read(reader, binary.BigEndian, &question.Class); err != nil {
		return question, err
	}
	return question, nil
}

// 解析DNS资源记录
func parseDNSRecord(reader *bytes.Reader) (DNSRecord, error) {
	var record DNSRecord
	name, err := parseDNSName(reader)
	if err != nil {
		return record, err
	}
	record.Name = name
	if err := binary.Read(reader, binary.BigEndian, &record.Type); err != nil {
		return record, err
	}
	if err := binary.Read(reader, binary.BigEndian, &record.Class); err != nil {
		return record, err
	}
	if err := binary.Read(reader, binary.BigEndian, &record.TTL); err != nil {
		return record, err
	}
	var dataLength uint16
	if err := binary.Read(reader, binary.BigEndian, &dataLength); err != nil {
		return record, err
	}
	record.Data = make([]byte, dataLength)
	if _, err := reader.Read(record.Data); err != nil {
		return record, err
	}
	return record, nil
}

// 解析DNS名称
func parseDNSName(reader *bytes.Reader) ([]byte, error) {
	var name []byte
	for {
		var length byte
		if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
			return nil, err
		}
		if length == 0 {
			break
		}
		if name == nil {
			name = make([]byte, 0)
		}
		label := make([]byte, length)
		if _, err := reader.Read(label); err != nil {
			return nil, err
		}
		name = append(name, label...)
		name = append(name, '.')
	}
	return name, nil
}

// 解析OPT记录数据
func parseOPTData(data []byte) {
	reader := bytes.NewReader(data)
	var udpPayloadSize uint16
	if err := binary.Read(reader, binary.BigEndian, &udpPayloadSize); err != nil {
		fmt.Println("Error reading UDP payload size:", err)
		return
	}
	fmt.Printf("UDP Payload Size: %d\n", udpPayloadSize)

	// 解析其他选项
	for reader.Len() > 0 {
		var optionCode uint16
		var optionLength uint16
		if err := binary.Read(reader, binary.BigEndian, &optionCode); err != nil {
			fmt.Println("Error reading option code:", err)
			return
		}
		if err := binary.Read(reader, binary.BigEndian, &optionLength); err != nil {
			fmt.Println("Error reading option length:", err)
			return
		}
		optionData := make([]byte, optionLength)
		if _, err := reader.Read(optionData); err != nil {
			fmt.Println("Error reading option data:", err)
			return
		}
		fmt.Printf("Option Code: %d, Option Data: %v\n", optionCode, optionData)
	}
}

func main() {
	// 示例：假设我们已经有一个DNS请求消息的二进制数据
	data := []byte{
		// DNS消息头
		0x00, 0x01, // ID
		0x00, 0x00, // Flags
		0x00, 0x01, // QdCount
		0x00, 0x00, // AnCount
		0x00, 0x00, // NsCount
		0x00, 0x01, // ArCount
		// 问题部分
		0x03, 'w', 'w', 'w', 0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00, // Name
		0x00, 0x01, // Type (A)
		0x00, 0x01, // Class (IN)
		// 额外记录部分 (OPT记录)
		0x00,       // Name (root)
		0x00, 0x29, // Type (OPT)
		0x10, 0x00, // UDP payload size (4096)
		0x00, 0x00, 0x00, 0x00, // Extended RCODE and Version
		0x00, 0x00, // RDLENGTH
	}

	parseDNSMessage(data)
}
