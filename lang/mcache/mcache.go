// Copyright 2021 ByteDance Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mcache

import (
	"sync"
)

const maxSize = 46

// index contains []byte which cap is 1<<index
var caches *MCache

func init() {
	caches = New()
}

type MCache [maxSize]sync.Pool

func Malloc(size int, capacity ...int) []byte {
	return caches.Malloc(size, capacity...)
}

func Free(buf []byte) {
	caches.Free(buf)
}

func New() *MCache {
	mc := new(MCache)
	for i := 0; i < maxSize; i++ {
		size := 1 << i
		mc[i].New = func() interface{} {
			buf := make([]byte, 0, size)
			return buf
		}
	}
	return mc
}

// Malloc supports one or two integer argument.
// The size specifies the length of the returned slice, which means len(ret) == size.
// A second integer argument may be provided to specify the minimum capacity, which means cap(ret) >= cap.
func (m *MCache) Malloc(size int, capacity ...int) []byte {
	if len(capacity) > 1 {
		panic("too many arguments to Malloc")
	}
	var c = size
	if len(capacity) > 0 && capacity[0] > size {
		c = capacity[0]
	}
	var ret = m[calcIndex(c)].Get().([]byte)
	ret = ret[:size]
	return ret
}

// Free should be called when the buf is no longer used.
func (m *MCache) Free(buf []byte) {
	size := cap(buf)
	if !isPowerOfTwo(size) {
		return
	}
	buf = buf[:0]
	m[bsr(size)].Put(buf)
}

// calculates which pool to get from
func calcIndex(size int) int {
	if size == 0 {
		return 0
	}
	if isPowerOfTwo(size) {
		return bsr(size)
	}
	return bsr(size) + 1
}
