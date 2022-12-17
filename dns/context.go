package dns

import "sync"

type Context struct {
	rw      *udpResponseWriter
	req     *Message
	handler Handler
	stats   Stats
}

var CtxPool = &sync.Pool{
	New: func() interface{} {
		ctx := new(Context)
		ctx.rw = new(udpResponseWriter)
		ctx.req = new(Message)
		ctx.req.Raw = make([]byte, 0, 1024)
		ctx.req.Domain = make([]byte, 0, 256)
		return ctx
	},
}
