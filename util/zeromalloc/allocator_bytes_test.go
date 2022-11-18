package zeromalloc

import (
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesAllocator(t *testing.T) {
	is := assert.New(t)
	unit, page, limit := 1024, 1024*4, 1024*1024
	al, err := NewPLocal(unit, page, limit)
	is.NoError(err)
	bal, err := NewBytesAllocator(al)
	is.NoError(err)

	//testdata := []byte("hello")
	procs := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	for p := 0; p < procs; p++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 1024*1024*1024; i++ {
				buf, err := bal.Alloc()
				is.NoError(err)
				is.Equal(unit, len(buf))
				bal.Free(buf)
			}
		}()
	}
	wg.Done()
}
