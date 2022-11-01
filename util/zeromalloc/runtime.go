package zeromalloc

import (
	_ "unsafe"
)

//go:linkname procPin runtime.procPin
func procPin() int

//go:linkname procUnpin runtime.procUnpin
func procUnpin() int
