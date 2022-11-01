package zeromalloc

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

//go:notinheap
type testStruct struct {
	A int
	B uintptr
}

type testHeapStruct struct {
	A int
	B uintptr
}

func toObject(p uintptr) *testStruct {
	return (*testStruct)(unsafe.Pointer(p))
}

type allocatorNewer func(unit int, page int, limit int) (Allocator, error)

type testcase struct {
	Name  string
	Newer allocatorNewer
}

func TestAllocator(t *testing.T) {
	is := assert.New(t)
	runtime.GOMAXPROCS(1)

	testcases := []testcase{
		{Name: "Unsafe", Newer: NewUnsafe},
		{Name: "Safe", Newer: NewSafe},
		{Name: "PLocal", Newer: NewPLocal},
	}
	page, limit := 1024*4, 1024*1024
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			a, err := tc.Newer(int(unsafe.Sizeof(testStruct{})), page, limit)
			is.NoError(err)
			ptr, err := a.Alloc()
			is.NoError(err)
			obj := toObject(ptr)
			obj.A = 1
			obj.B = 2
			runtime.GC()
			is.Equal(1, obj.A)
			is.Equal(uintptr(2), obj.B)
			a.Free(ptr)

			ptr, _ = a.Alloc()
			obj = toObject(ptr)
			runtime.GC()
			is.Equal(1, obj.A)
			is.Equal(uintptr(2), obj.B)
		})
	}
}
