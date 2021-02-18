package internal

import (
	"runtime"
	"sync/atomic"
)

type Pool struct {
	newFunc    func() any
	pool       []*_PoolElem
	Callers    [][]byte
	capacity   int32
	n          int32
	LogCallers bool
}

type _PoolElem struct {
	Value any
	Taken int32
}

func NewPool(
	capacity int32,
	newFunc func() any,
) *Pool {

	pool := &Pool{
		capacity: capacity,
		newFunc:  newFunc,
	}

	for i := int32(0); i < capacity; i++ {
		pool.pool = append(pool.pool, &_PoolElem{
			Value: newFunc(),
		})
	}
	pool.Callers = make([][]byte, capacity)

	return pool
}

func (p *Pool) Get() (value any, put func()) {
	idx := atomic.AddInt32(&p.n, 1) % p.capacity
	elem := p.pool[idx]
	if atomic.CompareAndSwapInt32(&elem.Taken, 0, 1) {
		if p.LogCallers {
			stack := make([]byte, 8*1024)
			runtime.Stack(stack, false)
			p.Callers[idx] = stack
		}
		value = elem.Value
		put = func() {
			if p.LogCallers {
				p.Callers[idx] = nil
			}
			atomic.StoreInt32(&elem.Taken, 0)
		}
	} else {
		value = p.newFunc()
		put = func() {}
	}
	return
}
