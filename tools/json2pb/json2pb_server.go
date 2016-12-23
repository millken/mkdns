package main

import (
	"encoding/json"
	"encoding/base64"
	"net/http"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/millken/mkdns/types"
)

func Json2Protobuf(c *gin.Context) {
	var err error
	records := types.Records{}

	content := c.DefaultPostForm("content", "")
	if content == "" {
		c.JSON(200, gin.H{"status": 401, "info": "field content empty"})
		return
	}
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		c.JSON(200, gin.H{"status": 405, "info": "base64 decode err"})
		return
	}
	err = json.Unmarshal(data, &records)
	if err != nil {
		log.Fatalln("json unmarshal error: ", err)
		c.JSON(200, gin.H{"status": 402, "info": "json decode err"})
		return
	}
	
	data, err = proto.Marshal(&records)
	if err != nil {
		log.Fatalln("protobuf marshal error: ", err)
		c.JSON(200, gin.H{"status": 403, "info": "proto encode err"})
		return
	}

	c.JSON(200, gin.H{"status": 200, "pb": base64.StdEncoding.EncodeToString(data)})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
    ApiServer := gin.Default()
    ApiServer.POST("/js2pb", Json2Protobuf)
	suser := &http.Server{
		Addr:           ":19180",
		Handler:        ApiServer,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 10,
	}
	suser.ListenAndServe()
}
