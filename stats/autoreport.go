package stats

import (
	"bytes"
	"expvar"
	"log"
	"net/http"

	"github.com/robfig/cron"
)

type AutoReportHandle struct {
	url      string
	schedule string
}

func NewAutoReport(url, schedule string) *AutoReportHandle {
	s := new(AutoReportHandle)
	s.url = url
	s.schedule = schedule
	return s
}

func (s *AutoReportHandle) report() {
	val := expvar.Get(statsVar)
	if val == nil || len(counters) == 0 {
		return
	}
	http.Post(s.url, "application/json; charset=utf-8", bytes.NewBuffer([]byte(val.String())))
	CleanCounter()
	log.Printf("[INFO] post data to :%s [%s]", s.url, val.String())
}

func (s *AutoReportHandle) Start() {
	c := cron.New()
	c.AddFunc("0 */1 * * * *", func() { s.report() })
	c.Start()
}
