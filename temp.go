package main

import (
	"container/ring"
	"strconv"
	"sync"
	"time"
)

type TempBuffer struct {
	buf *ring.Ring
	sync.Mutex
}

func (b *TempBuffer) Append(t Temp) {
	b.Lock()
	b.buf.Value = t
	b.buf = b.buf.Next()
	b.Unlock()
}

func (b *TempBuffer) Slice() []Temp {
	b.Lock()

	var vals []Temp
	b.buf.Do(func(f interface{}) {
		if f != nil {
			vals = append(vals, f.(Temp))
		}
	})
	b.Unlock()
	return vals
}

type Temp struct {
	TS    time.Time
	Value float32
}

func (t Temp) String() string {
	return "[" + t.TS.String() + "," + strconv.FormatFloat(float64(t.Value), 'f', 4, 32) + "]"
}
