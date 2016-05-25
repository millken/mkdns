package backends

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/millken/mkdns/types"
)

func init() {
	Register("file", NewFileBackend)
}

type FileBackend struct {
	zonefile string
}

func NewFileBackend(u *url.URL) (Backend, error) {
	zfile := u.Host + u.Path
	_, err := os.Open(zfile)
	if err != nil {
		return nil, fmt.Errorf("[%s]: %s", zfile, err)
	}
	fileBackend := &FileBackend{
		zonefile: zfile,
	}
	return fileBackend, nil
}

func (f *FileBackend) Load() error {
	zoneContent, err := ioutil.ReadFile(f.zonefile)
	if err != nil {
		return err
	}
	log.Printf("[INFO] loading zone: %s", f.zonefile)
	blists := string(zoneContent)
	blist := strings.Split(blists, "\n")
	reg := regexp.MustCompile(`\s+|\t+`)

	temp := make(map[string]*types.Zone)
	for _, b := range blist {
		b = reg.ReplaceAllString(b, " ")
		if len(b) < 1 {
			continue
		}

		alist := strings.Split(b, " ")
		zone := types.NewZone()
		zone.Name = alist[0]
		body, err := ioutil.ReadFile(alist[1])
		if err != nil {
			log.Printf("[ERROR] read zone file : %s, SKIPPED", err)
			continue
		}

		err = zone.ParseBody(body)

		if err != nil {
			log.Printf("[ERROR] parse zone body : %s", err)
			continue
		}
		temp[alist[0]] = zone

	}
	zonesLock.Lock()
	zones = temp
	zonesLock.Unlock()
	return nil
}

func (f *FileBackend) Watch() {
}
