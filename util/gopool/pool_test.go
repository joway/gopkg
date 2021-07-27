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

package gopool

import (
	"sync"
	"sync/atomic"
	"testing"
)

func testPanicFunc() {
	panic("test")
}

func TestPool(t *testing.T) {
	p := NewPool("test", 100, NewConfig())
	var n int32
	var wg sync.WaitGroup
	for i := 0; i < 2000; i++ {
		wg.Add(1)
		p.Go(func() {
			defer wg.Done()
			atomic.AddInt32(&n, 1)
		})
	}
	wg.Wait()
	if n != 2000 {
		t.Error(n)
	}
}

func TestPoolPanic(t *testing.T) {
	p := NewPool("test", 100, NewConfig())
	p.Go(testPanicFunc)
}
