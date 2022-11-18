package zeromalloc

import (
	"errors"
	"unsafe"
)

var (
	_          Allocator = (*allocator)(nil)
	OutOfLimit           = errors.New("out of allocator limit")
)

var _linkPtr uintptr

const (
	prtSizeUintPtr = unsafe.Sizeof(_linkPtr)
	ptrSize        = int(prtSizeUintPtr)
)

type allocator struct {
	unit  int // element size
	size  int // allocator size
	page  int // page size
	limit int // allocator limit
	head  uintptr
	cache []uintptr
}

func NewUnsafe(unit int, page int, limit int) (Allocator, error) {
	if page < unit {
		return nil, errors.New("invalid page size: page should >= unit")
	}
	if limit > 0 && page > limit {
		return nil, errors.New("invalid limit size: limit should >= page")
	}
	return &allocator{
		unit:  unit,
		size:  0,
		page:  page,
		limit: limit,
	}, nil
}

func (a *allocator) Alloc() (uintptr, error) {
	if a.head == 0 {
		block := ptrSize + a.unit
		blocknum := a.page / block
		if blocknum <= 0 {
			blocknum = 1
		}
		nextgrow := blocknum * block
		nextsize := a.size + nextgrow
		if a.limit > 0 && nextsize > a.limit {
			return 0, OutOfLimit
		}
		head, err := mmap(nextgrow)
		if err != nil {
			return 0, err
		}
		a.cache = append(a.cache, head)
		a.size = nextsize
		var nextptr *uintptr
		for i := 0; i < blocknum-1; i++ {
			nextptr = (*uintptr)(unsafe.Pointer(head + uintptr(i*block)))
			*nextptr = head + uintptr((i+1)*block)
		}
		nextptr = (*uintptr)(unsafe.Pointer(head + uintptr((blocknum-1)*block)))
		*nextptr = 0
		a.head = head
	}
	v := a.head
	a.head = *((*uintptr)(unsafe.Pointer(v)))
	return v + prtSizeUintPtr, nil
}

func (a *allocator) Free(p uintptr) {
	nextptr := (*uintptr)(unsafe.Pointer(p - prtSizeUintPtr))
	*nextptr = a.head
	a.head = p - prtSizeUintPtr
}

func (a *allocator) Close() error {
	return nil
}

func (a *allocator) Unit() int {
	return a.unit
}
