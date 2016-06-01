package backends

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/millken/mkdns/types"
)

func init() {
	Register("file", NewFileBackend)
}

type FileBackend struct {
	root string
}

func NewFileBackend(u *url.URL) (Backend, error) {
	root := u.Host + u.Path
	_, err := os.Open(root)
	if err != nil {
		return nil, fmt.Errorf("[%s]: %s", root, err)
	}
	fileBackend := &FileBackend{
		root: root,
	}
	return fileBackend, nil
}

func (f *FileBackend) walk(path string, fi os.FileInfo, err error) error {
	if fi.IsDir() == false {
		log.Printf("[DEBUG] path=%s", path)
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("[ERROR] read file [%s] err: %s", path, err)
		} else {
			dpb, err := types.DecodeByProtobuff(content)
			if err != nil {
				log.Printf("[ERROR] DecodeByProtobuf[%s] err: %s", path, err)
			} else {
				name := fi.Name()
				if dpb.Domain != "" {
					name = dpb.Domain
				}
				zonemap.Set(name, content)
			}
		}
	}
	return nil
}

func (f *FileBackend) Load() {
	log.Printf("[INFO] loading root : %s", f.root)
	filepath.Walk(f.root, f.walk)
}

func (f *FileBackend) Watch() {
}
