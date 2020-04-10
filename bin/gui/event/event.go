package event

import (
	"log"
	"sync"
)

// Consumer message Consumer
type Consumer func(e string, param ...interface{}) error

type consumerItem struct {
	Event string
	CB    Consumer
}

type manager struct {
	cbFuncs map[int]consumerItem
	mu      sync.Mutex
	id      int
}

var mgr manager

func init() {
	mgr.cbFuncs = make(map[int]consumerItem)
}

// RegisterConsumer register consumer,if e is empty, will recieve all event
func RegisterConsumer(e string, cb Consumer) int {
	if cb == nil {
		return 0
	}
	it := consumerItem{e, cb}
	mgr.mu.Lock()
	mgr.id++
	id := mgr.id
	mgr.cbFuncs[id] = it
	mgr.mu.Unlock()
	return id
}

// Unregister unregister
func Unregister(id int) {
	mgr.mu.Lock()
	delete(mgr.cbFuncs, id)
	mgr.mu.Unlock()
}

// Send send message
func Send(e string, param ...interface{}) {
	mgr.mu.Lock()
	cbs := make([]Consumer, 0, len(mgr.cbFuncs))
	for _, it := range mgr.cbFuncs {
		if it.Event != "" && it.Event != e {
			continue
		}
		cbs = append(cbs, it.CB)
	}
	mgr.mu.Unlock()
	for _, cb := range cbs {
		go func(e string, cb Consumer, param ...interface{}) {
			defer recover()
			err := cb(e, param...)
			if err != nil {
				log.Println("fail to process message:", e, cb, err)
			}
		}(e, cb, param...)
	}
}
