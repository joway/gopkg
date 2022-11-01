package zeromalloc

import (
	"runtime"
	_ "unsafe"
)

var _ Allocator = (*pallocator)(nil)

// p local allocator
type pallocator struct {
	allocators []Allocator // [pid]Allocator
}

func NewPLocal(unit int, page int, limit int) (Allocator, error) {
	procs := runtime.GOMAXPROCS(0)
	allocators := make([]Allocator, procs)
	for i := 0; i < procs; i++ {
		a, err := NewUnsafe(unit, page, limit)
		if err != nil {
			return nil, err
		}
		allocators[i] = a
	}

	return &pallocator{
		allocators: allocators,
	}, nil
}

func (a *pallocator) Alloc() (p uintptr, err error) {
	pid := procPin()
	p, err = a.allocators[pid].Alloc()
	procUnpin()
	return p, err
}

func (a *pallocator) Free(p uintptr) {
	pid := procPin()
	a.allocators[pid].Free(p)
	procUnpin()
}

func (a *pallocator) Close() (err error) {
	var cerr error
	for _, alctr := range a.allocators {
		cerr = alctr.Close()
		if err == nil && cerr != nil {
			err = cerr
		}
	}
	return err
}
