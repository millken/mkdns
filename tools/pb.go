package main

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RecordPb_Type int32

const (
	RecordPb_SOA RecordPb_Type = 0
	RecordPb_NS  RecordPb_Type = 1
	RecordPb_A   RecordPb_Type = 2
)

var RecordPb_Type_name = map[int32]string{
	0: "SOA",
	1: "NS",
	2: "A",
}
var RecordPb_Type_value = map[string]int32{
	"SOA": 0,
	"NS":  1,
	"A":   2,
}

func (x RecordPb_Type) String() string {
	return proto.EnumName(RecordPb_Type_name, int32(x))
}
func (RecordPb_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type RecordPb struct {
	Name  string            `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Ttl   int32             `protobuf:"varint,2,opt,name=ttl" json:"ttl,omitempty"`
	Type  RecordPb_Type     `protobuf:"varint,3,opt,name=type,enum=types.RecordPb_Type" json:"type,omitempty"`
	Value []*RecordPb_Value `protobuf:"bytes,4,rep,name=value" json:"value,omitempty"`
}

func (m *RecordPb) Reset()                    { *m = RecordPb{} }
func (m *RecordPb) String() string            { return proto.CompactTextString(m) }
func (*RecordPb) ProtoMessage()               {}
func (*RecordPb) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *RecordPb) GetValue() []*RecordPb_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

type RecordPb_Value struct {
	Record    []string            `protobuf:"bytes,1,rep,name=record" json:"record,omitempty"`
	View      string              `protobuf:"bytes,2,opt,name=view" json:"view,omitempty"`
	Weight    int32               `protobuf:"varint,3,opt,name=weight" json:"weight,omitempty"`
	Continent string              `protobuf:"bytes,4,opt,name=continent" json:"continent,omitempty"`
	Country   string              `protobuf:"bytes,5,opt,name=country" json:"country,omitempty"`
	Soa       *RecordPb_Value_SOA `protobuf:"bytes,6,opt,name=soa" json:"soa,omitempty"`
}

func (m *RecordPb_Value) Reset()                    { *m = RecordPb_Value{} }
func (m *RecordPb_Value) String() string            { return proto.CompactTextString(m) }
func (*RecordPb_Value) ProtoMessage()               {}
func (*RecordPb_Value) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *RecordPb_Value) GetSoa() *RecordPb_Value_SOA {
	if m != nil {
		return m.Soa
	}
	return nil
}

type RecordPb_Value_SOA struct {
	Mname   string `protobuf:"bytes,1,opt,name=mname" json:"mname,omitempty"`
	Nname   string `protobuf:"bytes,2,opt,name=nname" json:"nname,omitempty"`
	Serial  int32  `protobuf:"varint,3,opt,name=serial" json:"serial,omitempty"`
	Refresh int32  `protobuf:"varint,4,opt,name=refresh" json:"refresh,omitempty"`
	Retry   int32  `protobuf:"varint,5,opt,name=retry" json:"retry,omitempty"`
	Expire  int32  `protobuf:"varint,6,opt,name=expire" json:"expire,omitempty"`
	Minttl  int32  `protobuf:"varint,7,opt,name=minttl" json:"minttl,omitempty"`
}

func (m *RecordPb_Value_SOA) Reset()                    { *m = RecordPb_Value_SOA{} }
func (m *RecordPb_Value_SOA) String() string            { return proto.CompactTextString(m) }
func (*RecordPb_Value_SOA) ProtoMessage()               {}
func (*RecordPb_Value_SOA) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0, 0} }

type ZonePb struct {
	Name    string      `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Records []*RecordPb `protobuf:"bytes,2,rep,name=records" json:"records,omitempty"`
}

func (m *ZonePb) Reset()                    { *m = ZonePb{} }
func (m *ZonePb) String() string            { return proto.CompactTextString(m) }
func (*ZonePb) ProtoMessage()               {}
func (*ZonePb) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ZonePb) GetRecords() []*RecordPb {
	if m != nil {
		return m.Records
	}
	return nil
}

func init() {
	proto.RegisterType((*RecordPb)(nil), "types.RecordPb")
	proto.RegisterType((*RecordPb_Value)(nil), "types.RecordPb.Value")
	proto.RegisterType((*RecordPb_Value_SOA)(nil), "types.RecordPb.Value.SOA")
	proto.RegisterType((*ZonePb)(nil), "types.ZonePb")
	proto.RegisterEnum("types.RecordPb_Type", RecordPb_Type_name, RecordPb_Type_value)
}

var fileDescriptor0 = []byte{
	// 365 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x92, 0xcf, 0x4e, 0x2a, 0x31,
	0x14, 0x87, 0xef, 0xfc, 0xe9, 0x70, 0x39, 0x24, 0xf7, 0x4e, 0x1a, 0x34, 0x95, 0x18, 0x43, 0x58,
	0x61, 0x4c, 0x66, 0x81, 0x4f, 0xc0, 0xca, 0x9d, 0x9a, 0x62, 0x5c, 0xb8, 0x1b, 0xf0, 0x28, 0x93,
	0x40, 0x4b, 0x3a, 0x05, 0xc4, 0x07, 0xf0, 0x5d, 0x7c, 0x27, 0x1f, 0xc6, 0x9e, 0x96, 0x89, 0x86,
	0xe8, 0xae, 0xdf, 0xaf, 0xa7, 0x67, 0xbe, 0x69, 0x0f, 0xc0, 0xab, 0x56, 0x58, 0xac, 0x8c, 0xb6,
	0x9a, 0x33, 0xbb, 0x5b, 0x61, 0x3d, 0x78, 0x4b, 0xe1, 0xaf, 0xc4, 0x99, 0x36, 0x8f, 0xb7, 0x53,
	0xce, 0x21, 0x55, 0xe5, 0x12, 0x45, 0xd4, 0x8f, 0x86, 0x6d, 0xe9, 0xd7, 0x3c, 0x87, 0xc4, 0xda,
	0x85, 0x88, 0x5d, 0xc4, 0x24, 0x2d, 0xf9, 0x10, 0x52, 0x3a, 0x2b, 0x12, 0x17, 0xfd, 0x1b, 0x75,
	0x0b, 0xdf, 0xa8, 0x68, 0x9a, 0x14, 0x77, 0x0e, 0xa5, 0xaf, 0xe0, 0x17, 0xc0, 0x36, 0xe5, 0x62,
	0x8d, 0x22, 0xed, 0x27, 0xc3, 0xce, 0xe8, 0xe8, 0xb0, 0xf4, 0x9e, 0x36, 0x65, 0xa8, 0xe9, 0x7d,
	0xc4, 0xc0, 0x7c, 0xc0, 0x8f, 0x21, 0x33, 0xbe, 0xc4, 0x89, 0x24, 0x4e, 0x64, 0x4f, 0xa4, 0xb7,
	0xa9, 0x70, 0xeb, 0x5d, 0x9c, 0x1e, 0xad, 0xa9, 0x76, 0x8b, 0xd5, 0xf3, 0xdc, 0x7a, 0x1d, 0x26,
	0xf7, 0xc4, 0x4f, 0xa1, 0x3d, 0xd3, 0xca, 0x56, 0x0a, 0x95, 0x75, 0x9f, 0xa7, 0x03, 0x5f, 0x01,
	0x17, 0xd0, 0x9a, 0xe9, 0xb5, 0xb2, 0x66, 0x27, 0x98, 0xdf, 0x6b, 0xd0, 0x29, 0x27, 0xb5, 0x2e,
	0x45, 0xe6, 0xd2, 0xce, 0xe8, 0xe4, 0x47, 0xe1, 0x62, 0x72, 0x33, 0x96, 0x54, 0xd5, 0x7b, 0x8f,
	0x20, 0x71, 0xc0, 0xbb, 0xc0, 0x96, 0xdf, 0x2e, 0x2e, 0x00, 0xa5, 0xca, 0xa7, 0xc1, 0x37, 0x00,
	0x09, 0xd7, 0x68, 0xaa, 0x72, 0xd1, 0x08, 0x07, 0x22, 0x25, 0x83, 0x4f, 0x06, 0xeb, 0xb9, 0xd7,
	0x65, 0xb2, 0x41, 0xea, 0x63, 0xb0, 0x51, 0x65, 0x32, 0x00, 0xf5, 0xc1, 0x97, 0x55, 0x65, 0xd0,
	0xbb, 0xba, 0x3e, 0x81, 0x28, 0x5f, 0x56, 0x8a, 0x9e, 0xac, 0x15, 0xf2, 0x40, 0x83, 0x33, 0x48,
	0xe9, 0x65, 0x78, 0xcb, 0x2b, 0xe7, 0x7f, 0x78, 0x06, 0xf1, 0xf5, 0x24, 0x8f, 0x38, 0x83, 0x68,
	0x9c, 0xc7, 0x83, 0x2b, 0xc8, 0x1e, 0xdc, 0x74, 0xfc, 0x32, 0x05, 0xe7, 0x64, 0x47, 0x97, 0x50,
	0xbb, 0xbf, 0xa1, 0xb7, 0xfc, 0x7f, 0x70, 0x35, 0xb2, 0xd9, 0x9f, 0x66, 0x7e, 0xbe, 0x2e, 0x3f,
	0x03, 0x00, 0x00, 0xff, 0xff, 0xca, 0x2c, 0x98, 0xdf, 0x6d, 0x02, 0x00, 0x00,
}

func main() {
	z := ZonePb{}
	data := `{"name":"test.com","records":[{"name":"@","ttl":300,"type":"SOA","value":[{"soa":{"mname":"ns1.test1.com.","nname":"dns-admin.test1.com.","serial":305419896,"refresh":1193046,"retry":624485,"expire":4913,"minttl":389333}}]},{"name":"@","ttl":300,"type":"NS","value":[{"record":["ns1.test1.com.","ns2.test2.com"]}]},{"name":"@","ttl":300,"type":"A","value":[{"record":["1.47.46.2","1.2.3.3"]}]},{"name":"view","ttl":300,"type":"A","state":1,"value":[{"record":["1.47.46.2","1.2.3.3"],"view":"any"},{"record":["1.2.3.4","1.2.4.5"],"view":"dx"}]},{"name":"weight","ttl":300,"type":"A","state":2,"value":[{"record":["1.47.46.2","1.2.3.3"],"weight":3},{"record":["1.2.3.4","1.2.4.5"],"weight":7},{"record":["7.7.7.7"],"weight":10}]},{"name":"geo","ttl":300,"type":"A","state":4,"value":[{"record":["7.7.7.7"]},{"record":["1.7.6.2"],"continent":"asia"},{"record":["1.7.6.5"],"continent":"asia","country":"cn"},{"record":["1.7.6.6"],"country":"cn"},{"record":["1.2.3.4","1.2.4.5"],"country":"kr"},{"record":["1.1.1.1","1.2.2.3","1.1.1.2"],"continent":"north-america"},{"record":["1.1.1.3"],"country":"us"}]}]}`
	if err := proto.Unmarshal([]byte(data), &z); err != nil {
		fmt.Println(err)
	}
}
