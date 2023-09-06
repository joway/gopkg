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

package runtimex

import (
	"unsafe"

	"github.com/modern-go/reflect2"
)

//go:linkname Fastrand runtime.fastrand
func Fastrand() uint32

func getg() uintptr

type puintptr uintptr
type guintptr uintptr
type muintptr uintptr

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
	pidField         = pType.FieldByName("id")
	pstatusField     = pType.FieldByName("status")
	psysmontickField = pType.FieldByName("sysmontick")
	ppreemptField    = pType.FieldByName("preempt")
	prunqheadField   = pType.FieldByName("runqhead")
	prunqtailField   = pType.FieldByName("runqtail")
	//prunqField       = pType.FieldByName("runq")
)

var _ G = (*_g)(nil)
var _ M = (*_m)(nil)
var _ P = (*_p)(nil)

type G interface {
	M() M
	Id() *int64
	Preempt() *bool
}

type M interface {
	P() P
	Id() *int64
	PreemptOff() *string
}

type P interface {
	Id() *int32
	Status() *uint32
	Sysmontick() *Sysmontick
	Preempt() *bool
	RunqSize() uint32
}

type Sysmontick struct {
	Schedtick   uint32
	Schedwhen   int64
	Syscalltick uint32
	Syscallwhen int64
}

func GetG() G {
	gp := getg()
	g := _g{
		ptr: gp,
		m:   newM(gp + mField.Offset()),
	}
	return g
}

type _g struct {
	ptr uintptr
	m   M
}

func (g _g) M() M {
	return g.m
}

func (g _g) Id() *int64 {
	gid := (*int64)(unsafe.Pointer(g.ptr + gidField.Offset()))
	return gid
}

func (g _g) Preempt() *bool {
	preempt := (*bool)(unsafe.Pointer(g.ptr + gpreemptField.Offset()))
	return preempt
}

func newM(mp uintptr) _m {
	m := *((**_m)(unsafe.Pointer(mp)))
	mptr := uintptr(unsafe.Pointer(m))
	pp := uintptr(unsafe.Pointer(m)) + pField.Offset()
	p := newP(pp)
	return _m{ptr: mptr, p: p}
}

type _m struct {
	ptr uintptr
	p   P
}

func (m _m) P() P {
	return m.p
}

func (m _m) Id() *int64 {
	id := (*int64)(unsafe.Pointer(m.ptr + midField.Offset()))
	return id
}
func (m _m) PreemptOff() *string {
	preemptoff := (*string)(unsafe.Pointer(m.ptr + mpreemptoffField.Offset()))
	return preemptoff
}

func newP(pp uintptr) _p {
	p := *((**_p)(unsafe.Pointer(pp)))
	pptr := uintptr(unsafe.Pointer(p))
	return _p{ptr: pptr}
}

type _p struct {
	ptr uintptr
}

func (p _p) Id() *int32 {
	id := (*int32)(unsafe.Pointer(p.ptr + pidField.Offset()))
	return id
}

func (p _p) Status() *uint32 {
	status := (*uint32)(unsafe.Pointer(p.ptr + pstatusField.Offset()))
	return status
}

func (p _p) Sysmontick() *Sysmontick {
	st := (*Sysmontick)(unsafe.Pointer(p.ptr + psysmontickField.Offset()))
	return st
}

func (p _p) Preempt() *bool {
	preempt := (*bool)(unsafe.Pointer(p.ptr + ppreemptField.Offset()))
	return preempt
}

func (p _p) RunqSize() uint32 {
	prunqhead := *(*uint32)(unsafe.Pointer(p.ptr + prunqheadField.Offset()))
	prunqtail := *(*uint32)(unsafe.Pointer(p.ptr + prunqtailField.Offset()))
	qsize := prunqtail - prunqhead
	return qsize
}
