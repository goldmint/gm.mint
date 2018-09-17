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

	// time
	if n.useTime {
		ts := uint64(time.Now().UTC().Unix() * 1000)
		if ts <= n.nonce {
			n.nonce++
		} else {
			n.nonce = ts
		}
		return n.nonce
	}

	// just nonce
	ret := n.nonce
	n.nonce++
	return ret
}
