package zeromalloc

import (
	"runtime"
	"sync/atomic"
)

type Allocator interface {
	Alloc() (uintptr, error)
	Free(p uintptr)
	Close() error
	Unit() int
}

type BytesAllocator interface {
	Alloc() (b []byte, err error)
	Free(b []byte)
	Close() error
}

func lock(mu *uint32) {
	for !atomic.CompareAndSwapUint32(mu, 0, 1) {
		runtime.Gosched()
	}
}

func unlock(mu *uint32) {
	for !atomic.CompareAndSwapUint32(mu, 1, 0) {
		runtime.Gosched()
	}
}
