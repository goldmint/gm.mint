package nonce

import (
	"sync"
	"time"
)

// New nonce holder
func New(val uint64, useTime bool) *Nonce {
	return &Nonce{
		nonce:   val,
		lock:    &sync.Mutex{},
		useTime: useTime,
	}
}

// ---

// Nonce data
type Nonce struct {
	nonce   uint64
	lock    *sync.Mutex
	useTime bool
}

// Next nonce
func (n *Nonce) Next() uint64 {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.useTime {
		ts := uint64(time.Now().UTC().UnixNano())
		if ts <= n.nonce {
			n.nonce++
		} else {
			n.nonce = ts
		}
		return n.nonce
	}

	n.nonce++
	return n.nonce
}
