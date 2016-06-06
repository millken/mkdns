package backends

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/millken/mkdns/types"
)

func init() {
	Register("file", NewFileBackend)
}

type FileBackend struct {
	root     string
	fsnotify *fsnotify.Watcher
	events   chan fsnotify.Event
}

func NewFileBackend(u *url.URL) (Backend, error) {
	root := u.Host + u.Path
	_, err := os.Open(root)
	if err != nil {
		return nil, fmt.Errorf("[%s]: %s", root, err)
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fileBackend := &FileBackend{
		root:     root,
		fsnotify: watcher,
		events:   make(chan fsnotify.Event),
	}
	return fileBackend, nil
}

func (f *FileBackend) read(path string) {
	var name string
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("[ERROR] read file [%s] err: %s", path, err)
	} else {
		dpb, err := types.DecodeByProtobuff(content)
		if err != nil {
			log.Printf("[ERROR] DecodeByProtobuf[%s] err: %s", path, err)
		} else {
			s, err := os.Stat(path)
			if err != nil {
				log.Printf("[ERROR] stat file %s err: %s", path, err)
			} else {
				name = s.Name()
			}
			if dpb.Domain != "" {
				name = dpb.Domain
			}
			zonemap.Set(name, content)
			zonecache.Del(name)
			log.Printf("[DEBUG] path=%s, name=%s, dpb = %+v", path, name, dpb)
		}
	}

}

func (f *FileBackend) walk(path string, fi os.FileInfo, err error) error {
	if fi.IsDir() == false {
		f.read(path)
	} else {
		if err := f.fsnotify.Add(path); err != nil {
			log.Printf("[WARN] file watch : %s", err)
		}

	}
	return nil
}

func (f *FileBackend) Load() {
	log.Printf("[INFO] loading root : %s", f.root)
	go f.watch()

	start := time.Now()
	filepath.Walk(f.root, f.walk)
	end := time.Now()
	log.Printf("[INFO] loaded root : %s cost time : %v\n", f.root, end.Sub(start))
}

func (f *FileBackend) watch() {
	defer f.fsnotify.Close()
	for {
		select {

		case e := <-f.fsnotify.Events:
			s, err := os.Stat(e.Name)
			if err == nil && s != nil {
				if s.IsDir() {
					if e.Op&fsnotify.Create != 0 {
						log.Printf("[DEBUG] add dir %s to watcher", e.Name)
						f.fsnotify.Add(e.Name)
					}

				} else {
					f.read(e.Name)
				}
			}
			if e.Op&fsnotify.Remove != 0 {
				if filepath.Ext(e.Name) == "" {
					log.Printf("[DEBUG] remove dir %s from watcher", e.Name)
					f.fsnotify.Remove(e.Name)
				} else {
					fname := filepath.Base(e.Name)
					log.Printf("[INFO] domain config removed : %s", fname)
					zonemap.Remove(fname)
					zonecache.Del(fname)
				}
			}

		case e := <-f.fsnotify.Errors:
			log.Printf("[ERROR] notify err: %s", e)
		}
	}
}
