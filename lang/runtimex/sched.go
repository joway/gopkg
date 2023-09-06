package runtimex

import (
	"log"
	"unsafe"

	"github.com/modern-go/reflect2"
)

var (
	// types
	gType = reflect2.TypeByName("runtime.g").(reflect2.StructType)
	mType = reflect2.TypeByName("runtime.m").(reflect2.StructType)
	pType = reflect2.TypeByName("runtime.p").(reflect2.StructType)
	// fields
	mField           = gType.FieldByName("m")
	pField           = mType.FieldByName("p")
	gidField         = gType.FieldByName("goid")
	gpreemptField    = gType.FieldByName("preempt")
	midField         = mType.FieldByName("id")
	mpreemptoffField = mType.FieldByName("preemptoff")
	ppreemptField    = pType.FieldByName("preempt")
	prunqheadField   = pType.FieldByName("runqhead")
	prunqtailField   = pType.FieldByName("runqtail")
	//prunqField       = pType.FieldByName("runq")
)

func MPreemptOff() {
	g := getg()
	m := *((**_m)(unsafe.Pointer(g + mField.Offset())))
	preemptoff := (*string)(unsafe.Pointer((uintptr)(unsafe.Pointer(m)) + mpreemptoffField.Offset()))
	if *preemptoff == "" {
		*preemptoff = "holding"
	}
}

func MPreemptOn() {
	g := getg()
	m := *((**_m)(unsafe.Pointer(g + mField.Offset())))
	preemptoff := (*string)(unsafe.Pointer((uintptr)(unsafe.Pointer(m)) + mpreemptoffField.Offset()))
	if *preemptoff != "" {
		*preemptoff = ""
	}
}

func schedlog() {
	g := getg()
	m := *((**_m)(unsafe.Pointer(g + mField.Offset())))
	p := *((**_p)(unsafe.Pointer(uintptr(unsafe.Pointer(m)) + pField.Offset())))
	gid := *((*int64)(unsafe.Pointer(g + gidField.Offset())))
	gpreempt := *(*bool)(unsafe.Pointer(g + gpreemptField.Offset()))
	mid := *((*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(m)) + midField.Offset())))
	mpreemptoff := *(*string)(unsafe.Pointer((uintptr)(unsafe.Pointer(m)) + mpreemptoffField.Offset()))
	pid := p.id
	pstatus := p.status
	ppreempt := *(*bool)(unsafe.Pointer((uintptr)(unsafe.Pointer(p)) + ppreemptField.Offset()))
	prunqhead := *(*uint32)(unsafe.Pointer((uintptr)(unsafe.Pointer(p)) + prunqheadField.Offset()))
	prunqtail := *(*uint32)(unsafe.Pointer((uintptr)(unsafe.Pointer(p)) + prunqtailField.Offset()))
	//prunq := *(*[256]guintptr)(unsafe.Pointer((uintptr)(unsafe.Pointer(p)) + prunqField.Offset()))
	qsize := prunqtail - prunqhead
	log.Printf(
		"[G] gid=%d,gpreempt=%v | "+
			"[M] mid=%d,mpreemptoff=%s | "+
			"[P] pid=%d,pstatus=%d,ppreempt=%v,qsize=%d",
		gid, gpreempt,
		mid, mpreemptoff,
		pid, pstatus, ppreempt, qsize,
	)
}
