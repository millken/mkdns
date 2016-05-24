package backend

import "log"

type Backend interface {
	Get(name string) error
	Watch()
}

var Backends = map[string]func(string) (Backend, error){}

func Register(name string, backend func(string) (Backend, error)) {
	if _, dup := Backends[name]; dup {
		log.Fatal("duplicate backend", name)
	}
	Backends[name] = backend
}
