package stats

import (
	"encoding/json"
	"expvar"
	"net/http"
	"sync"
)

const (
	statsVar = "stats.counter"
)

type statsHandler struct {
	stats     map[string]uint64
	statsLock sync.Mutex
}

func NewStats() *statsHandler {
	s := new(statsHandler)
	s.stats = make(map[string]uint64)
	return s
}

func writeJsonResponse(w http.ResponseWriter, obj interface{}, err error) error {
	if err == nil {
		encoder := json.NewEncoder(w)
		encoder.Encode(obj)
		return nil
	}
	return err
}

func (s *statsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate,no-cache,no-store")

	val := expvar.Get(statsVar)
	if val == nil {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("No metrics."))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(val.String()))
}
