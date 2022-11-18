package zeromalloc

import (
	"runtime"
	"sync"
	"testing"
	"unsafe"
)

func BenchmarkSequentialAllocator(b *testing.B) {
	page, limit := 1024*4, 1024*1024
	benchcases := []testcase{
		{Name: "Unsafe", Newer: NewUnsafe},
		{Name: "Safe", Newer: NewSafe},
		{Name: "PLocal", Newer: NewPLocal},
	}
	for _, bc := range benchcases {
		b.Run(bc.Name, func(b *testing.B) {
			al, _ := bc.Newer(int(unsafe.Sizeof(testStruct{})), page, limit)
			for i := 0; i < b.N; i++ {
				ptr, _ := al.Alloc()
				obj := toObject(ptr)
				obj.A = i
				obj.B = uintptr(i)
				al.Free(ptr)
			}
		})
	}
}

func BenchmarkConcurrentAllocator(b *testing.B) {
	b.Run("SyncPool", func(b *testing.B) {
		pool := sync.Pool{New: func() interface{} {
			return &testHeapStruct{}
		}}
		b.RunParallel(func(pb *testing.PB) {
			var sum int
			for pb.Next() {
				sum++
				p := pool.Get()
				obj := p.(*testHeapStruct)
				obj.A = sum
				obj.B = uintptr(sum)
				pool.Put(obj)
			}
		})
	})
	b.Run("PLocal", func(b *testing.B) {
		page, limit := 1024*4, 1024*1024
		b.RunParallel(func(pb *testing.PB) {
			al, _ := NewPLocal(int(unsafe.Sizeof(testStruct{})), page, limit)
			var sum int
			for pb.Next() {
				sum++
				ptr, _ := al.Alloc()
				obj := toObject(ptr)
				obj.A = sum
				obj.B = uintptr(sum)
				al.Free(ptr)
			}
		})
	})
}

func BenchmarkBytesAllocator(b *testing.B) {
	unit := 1024
	b.Run("SyncPool", func(b *testing.B) {
		pool := sync.Pool{New: func() interface{} {
			return make([]byte, 1024)
		}}
		b.RunParallel(func(pb *testing.PB) {
			var sum int
			for pb.Next() {
				if sum%100 == 0 {
					runtime.GC()
				}
				sum++
				p := pool.Get()
				obj := p.([]byte)
				copy(obj, []byte("hello"))
				pool.Put(obj)
			}
		})
	})
	b.Run("PLocal", func(b *testing.B) {
		page, limit := 1024*4, 1024*1024
		b.RunParallel(func(pb *testing.PB) {
			al, _ := NewPLocal(unit, page, limit)
			bal, _ := NewBytesAllocator(al)
			var sum int
			for pb.Next() {
				if sum%100 == 0 {
					runtime.GC()
				}
				sum++
				obj, _ := bal.Alloc()
				copy(obj, []byte("hello"))
				bal.Free(obj)
			}
		})
	})
}
