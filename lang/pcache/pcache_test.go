package pcache

import (
	"fmt"
	"testing"

	"github.com/bytedance/gopkg/lang/mcache"
)

func BenchmarkCache(b *testing.B) {
	type benchcase struct {
		Name   string
		Malloc func(size int, capacity ...int) []byte
		Free   func(buf []byte)
	}
	benchcases := []benchcase{
		{Name: "MCache", Malloc: mcache.Malloc, Free: mcache.Free},
		{Name: "PCache", Malloc: Malloc, Free: Free},
	}
	sizecases := []int{256, 512, 1024, 4096, 10240, 102400}
	for _, cs := range benchcases {
		b.Run(cs.Name, func(b *testing.B) {
			for _, size := range sizecases {
				b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
					b.RunParallel(func(pb *testing.PB) {
						origin := make([]byte, size)
						for i := 0; i < size; i++ {
							origin[i] = 'a' + byte(i%26)
						}
						b.ResetTimer()
						b.ReportAllocs()
						for pb.Next() {
							buf := cs.Malloc(size)
							copy(buf, origin)
							cs.Free(buf)
						}
					})
				})
			}
		})
	}
}
