package backends

import (
	"crypto/rc4"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/millken/mkdns/types"
	"github.com/mreiferson/go-httpclient"
)

func init() {
	Register("http", NewHttpBackend)
}

// HttpBackend postgres as backend store
type HttpBackend struct {
	url       *url.URL
	lastUtime string
	requestId int
	appId     string
	appKey    string
	host      string
	version   string
}

type StMax struct {
	Requestid int `json:"requestid"`
	Data      struct {
		Maxid int    `json:"maxid"`
		Total int    `json:"total"`
		Utime string `json:"utime"`
	} `json:"data"`
	Status int `json:"status"`
}

type StEvent struct {
	Data struct {
		Records []struct {
			Domain string      `json:"Domain"`
			Act    string      `json:"Act"`
			Utime  string      `json:"Utime"`
			Value  interface{} `json:"Value"`
		} `json:"records"`
	} `json:"data"`
	Requestid int `json:"requestid"`
	Status    int `json:"status"`
}

type StList struct {
	Data struct {
		Records []struct {
			ID     int    `json:"Id"`
			Domain string `json:"Domain"`
			Value  string `json:"Value"`
		} `json:"records"`
	} `json:"data"`
	Requestid int `json:"requestid"`
	Status    int `json:"status"`
}

func NewHttpBackend(u *url.URL) (Backend, error) {
	password, _ := u.User.Password()
	backend := &HttpBackend{
		url:       u,
		appId:     u.User.Username(),
		appKey:    password,
		host:      u.Host,
		requestId: -1,
		version:   "1.0",
	}
	return backend, nil
}

func (b *HttpBackend) httpGet(url string) (body []byte, err error) {
	transport := &httpclient.Transport{
		ConnectTimeout:        1 * time.Second,
		RequestTimeout:        10 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
	}
	defer transport.Close()

	client := &http.Client{Transport: transport}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func (b *HttpBackend) getMax() (res StMax, err error) {
	orld := fmt.Sprintf("v=%s&action=dns.getMax&requestid=%d", b.version, b.requestId+1)
	appen := b.encode(orld)

	url := fmt.Sprintf("http://%s/?appid=%s&appen=%s", b.host, b.appId, appen)
	body, err := b.httpGet(url)
	if err != nil {
		return
	}
	log.Printf("[DEBUG] http get %s: %s", url, body)
	err = json.Unmarshal(body, &res)
	return
}

func (b *HttpBackend) encode(param string) string {
	key := []byte(b.appKey)
	cipher, _ := rc4.NewCipher(key)
	data := make([]byte, len(param))
	cipher.XORKeyStream(data, []byte(param))
	appen := base64.StdEncoding.EncodeToString(data)
	return appen

}

func (b *HttpBackend) getList(maxid, limit, offset int) (res StList, err error) {
	orld := fmt.Sprintf("v=%s&action=dns.getList&maxid=%d&limit=%d&offset=%d&requestid=%d", b.version, maxid, limit, offset, b.requestId+1)
	appen := b.encode(orld)

	url := fmt.Sprintf("http://%s/?appid=%s&appen=%s", b.host, b.appId, appen)
	body, err := b.httpGet(url)
	if err != nil {
		return
	}
	log.Printf("[DEBUG] http get %s: %s", url, body)
	err = json.Unmarshal(body, &res)
	return
}

func (b *HttpBackend) getEvent() (res StEvent, err error) {
	orld := fmt.Sprintf("v=%s&action=dns.getLastEvent&utime=%s&requestid=%d", b.version, b.lastUtime, b.requestId+1)
	appen := b.encode(orld)

	url := fmt.Sprintf("http://%s/?appid=%s&appen=%s", b.host, b.appId, appen)
	body, err := b.httpGet(url)
	if err != nil {
		return
	}
	//log.Printf("[DEBUG] http get %s: %s", url, body)
	err = json.Unmarshal(body, &res)
	return
}

// Load load config
func (b *HttpBackend) Load() {

	log.Printf("[INFO] loading http backend: %s", b.url.String())
	res, err := b.getMax()
	if err != nil {
		log.Printf("[ERROR] getMax() : %s", err)
		return
	}

	if res.Status != 200 {
		log.Printf("[ERROR] get code : %d", res.Status)
		return
	}
	b.requestId = res.Requestid
	b.lastUtime = res.Data.Utime

	log.Printf("[INFO] %+v", res)
	offset := 0
	limit := 100
	for offset < res.Data.Total {
		lists, err := b.getList(res.Data.Maxid, limit, offset)
		if err != nil {
			log.Printf("[ERROR] getList : %s", err)
			break
		}
		b.requestId = lists.Requestid

		offset = offset + len(lists.Data.Records)
		for _, l := range lists.Data.Records {
			data, err := base64.StdEncoding.DecodeString(l.Value)
			if err != nil {
				log.Printf("[ERROR]  err: %s", err)
				continue
			}
			dpb, err := types.DecodeByProtobuff(data)
			if err != nil {
				log.Printf("[ERROR] DecodeByProtobuf err: %s", err)
				continue
			}

			if dpb.Domain != "" {
				l.Domain = dpb.Domain
			}
			zonemap.Set(l.Domain, data)
			zonecache.Del(l.Domain)
			log.Printf("[DEBUG] ID: %d, domain: %s, dpb=%+v\n", l.ID, l.Domain, dpb)

		}
	}
	go b.watch()

}

func (b *HttpBackend) watch() {
	var isLocked bool
	ticker := time.NewTicker(time.Duration(3) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if isLocked {
				continue
			}
			isLocked = false
			res, err := b.getEvent()
			if err != nil {
				log.Printf("[ERROR] getEvent() : %s", err)
				continue
			}
			b.requestId = res.Requestid

			if res.Status != 200 {
				log.Printf("[ERROR] get code : %d", res.Status)
				continue
			}
			for _, l := range res.Data.Records {
				if l.Act == "delete" {
					log.Printf("[INFO] domain config removed : %s", l.Domain)
					zonemap.Remove(l.Domain)
					zonecache.Del(l.Domain)
					continue
				}
				if l.Value == nil {
					continue
				}
				data, err := base64.StdEncoding.DecodeString(l.Value.(string))
				if err != nil {
					log.Printf("[ERROR]  err: %s", err)
					continue
				}

				dpb, err := types.DecodeByProtobuff(data)
				if err != nil {
					log.Printf("[ERROR] DecodeByProtobuf err: %s", err)
					continue
				}
				if dpb.Domain != "" {
					l.Domain = dpb.Domain
				}
				zonemap.Set(l.Domain, data)
				zonecache.Del(l.Domain)
				b.lastUtime = l.Utime
				log.Printf("[INFO] domain config added : %s", l.Domain)

			}
		}
	}
}
