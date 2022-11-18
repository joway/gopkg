package zeromalloc

import (
	"reflect"
	"unsafe"
)

var _ BytesAllocator = (*ballocator)(nil)

type ballocator struct {
	Allocator
}

func NewBytesAllocator(a Allocator) (BytesAllocator, error) {
	return &ballocator{Allocator: a}, nil
}

func (a *ballocator) Alloc() (b []byte, err error) {
	var p uintptr
	p, err = a.Allocator.Alloc()
	if err != nil {
		return nil, err
	}
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh.Data = p
	sh.Len = a.Allocator.Unit()
	sh.Cap = a.Allocator.Unit()
	return b, nil
}

func (a *ballocator) Free(b []byte) {
	p := (*reflect.SliceHeader)(unsafe.Pointer(&b)).Data
	a.Allocator.Free(p)
}
