package types

import "github.com/golang/protobuf/proto"

func DecodeByProtobuff(data []byte) (zpb ZonePb, err error) {
	zpb = ZonePb{}
	err = proto.Unmarshal(data, &zpb)
	return
}
