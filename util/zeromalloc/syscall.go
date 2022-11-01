package zeromalloc

import (
	"syscall"
)

func mmap(size int) (uintptr, error) {
	fd := -1
	data, _, errno := syscall.Syscall6(
		syscall.SYS_MMAP,
		0, // address
		uintptr(size),
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_ANON|syscall.MAP_PRIVATE,
		uintptr(fd), // No file descriptor
		0,           // offset
	)
	if errno != 0 {
		return 0, syscall.Errno(errno)
	}
	return data, nil
}
