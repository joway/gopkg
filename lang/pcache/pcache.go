package pcache

import (
	"sync"

	"github.com/bytedance/gopkg/lang/mcache"
)

var defaultPCache = new(pcache)

func Malloc(size int, capacity ...int) []byte {
	return defaultPCache.Malloc(size, capacity...)
}

func Free(buf []byte) {
	defaultPCache.Free(buf)
}

type pcache struct {
	local sync.Pool
}

func (p *pcache) localPool() *mcache.MCache {
	x := p.local.Get()
	if x == nil {
		x = mcache.New()
	}
	mp := x.(*mcache.MCache)
	p.local.Put(x)
	return mp
}

func (p *pcache) Malloc(size int, capacity ...int) []byte {
	return p.localPool().Malloc(size, capacity...)
}

func (p *pcache) Free(buf []byte) {
	p.localPool().Free(buf)
}

type memPool struct {
	*mcache.MCache
}

func (p *memPool) Malloc(size int, capacity ...int) []byte {
	return p.Malloc(size, capacity...)
}

func (p *memPool) Free(buf []byte) {
	p.Free(buf)
}
