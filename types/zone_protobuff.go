package types

import "github.com/golang/protobuf/proto"

func DecodeByProtobuff(data []byte) (r Records, err error) {
	r = Records{}
	err = proto.Unmarshal(data, &r)
	return
}
