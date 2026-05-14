package main

import (
	. "net/http"
	"sync/atomic"
)

type hitsCounter struct {
	stats   *ApiStats
	handler Handler
}

func (c hitsCounter) ServeHTTP(w ResponseWriter, r *Request) {
	c.stats.hits.Add(1)
	c.handler.ServeHTTP(w, r)
}

type ApiStats struct {
	hits atomic.Int32
}

func (a *ApiStats) Reset() {
	a.hits.Store(0)
}

func (a *ApiStats) HitsCounter(handler Handler) Handler {
	return hitsCounter{
		stats:   a,
		handler: handler,
	}
}
