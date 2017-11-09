package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/millken/mkdns/types"
)

func main() {
	var err error
	records := types.Records{}
	in := flag.String("in", "test.pb", "json file path")
	out := flag.String("out", "test.json", "pb file path")
	flag.Parse()
	content, err := ioutil.ReadFile(*in)
	if err != nil {
		log.Fatalln("read json file err:", err)
	}
	if err := proto.Unmarshal(content, &records); err != nil {
		log.Fatalln("proto unmarshal error: ", err)
	}
	log.Printf("%+v", records)
	data, err := json.Marshal(records)
	if err != nil {
		log.Fatal("json marshal error: ", err)
	}
	err = ioutil.WriteFile(*out, data, 0666)
	if err != nil {
		log.Fatal("write file error: ", err)
	}
	log.Println("done")
}
