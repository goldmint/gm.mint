package nonce

import (
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	tasks := 4
	gens := 10000

	type TaskData struct {
		Res []uint64
	}

	// mode 1
	non := New(0, false)
	dat := []TaskData{}
	for i := 0; i < tasks; i++ {
		dat = append(dat, TaskData{
			Res: make([]uint64, 0, gens),
		})
	}
	wg := sync.WaitGroup{}
	for i := 0; i < tasks; i++ {
		wg.Add(1)
		go func(d *TaskData) {
			for g := 0; g < gens; g++ {
				d.Res = append(d.Res, non.Next())
			}
			wg.Done()
		}(&dat[i])
	}
	wg.Wait()

	// check duplicates
	dups := make(map[uint64]struct{})
	for i := 0; i < tasks; i++ {
		for _, v := range dat[i].Res {
			if _, has := dups[v]; has {
				t.Fatal("Has duplicate")
			} else {
				dups[v] = struct{}{}
			}
		}
	}

	t.Log("Mode 1: OK, there aren't duplicates in", len(dups), "items")

	// ---

	// mode 2
	non = New(0, true)
	dat = []TaskData{}
	for i := 0; i < tasks; i++ {
		dat = append(dat, TaskData{
			Res: make([]uint64, 0, gens),
		})
	}
	wg = sync.WaitGroup{}
	for i := 0; i < tasks; i++ {
		wg.Add(1)
		go func(d *TaskData) {
			for g := 0; g < gens; g++ {
				d.Res = append(d.Res, non.Next())
			}
			wg.Done()
		}(&dat[i])
	}
	wg.Wait()

	// check duplicates
	dups = make(map[uint64]struct{})
	for i := 0; i < tasks; i++ {
		for _, v := range dat[i].Res {
			if _, has := dups[v]; has {
				t.Fatal("Has duplicate")
			} else {
				dups[v] = struct{}{}
			}
		}
	}

	t.Log("Mode 2: OK, there aren't duplicates in", len(dups), "items")
}
