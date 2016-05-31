package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/millken/mkdns/types"
)

func main() {
	var err error
	zonepb := types.ZonePb{}
	in := flag.String("in", "test.json", "json file path")
	out := flag.String("out", "test.pb", "pb file path")
	flag.Parse()
	content, err := ioutil.ReadFile(*in)
	if err != nil {
		log.Fatalln("read json file err:", err)
	}
	//fmt.Printf("json body = %s", content)
	/*
		err = json.Unmarshal(content, &zonepb)
		if err != nil {
			log.Fatalln("json unmarshal error: ", err)
		}
	*/
	jsonpb.UnmarshalString(string(content[:]), &zonepb)
	log.Printf("%+v", zonepb)
	data, err := proto.Marshal(&zonepb)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	err = ioutil.WriteFile(*out, data, 0666)
	if err != nil {
		log.Fatal("write file error: ", err)
	}
	log.Println("done")
}
