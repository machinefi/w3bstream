package safe

import "sync/atomic"

type Bool struct{ int32 }

var (
	fJSON  = []byte("false")
	tJSON  = []byte("true")
	tInt32 = int32(1)
)

func NewBool() *Bool { return &Bool{} }

func NewBoolWithVal(v bool) *Bool {
	vb := NewBool()
	vb.Set(v)
	return vb
}

func (b *Bool) Clone() *Bool { return NewBoolWithVal(b.Val()) }

func (b *Bool) Val() bool { return atomic.LoadInt32(&b.int32) == 1 }

func (b *Bool) CAS(pv, nv bool) (swapped bool) {
	var prev, curr int32
	if pv {
		prev = tInt32
	}
	if nv {
		curr = tInt32
	}
	return atomic.CompareAndSwapInt32(&b.int32, prev, curr)
}

func (b *Bool) Set(v bool) (old bool) {
	if v {
		old = atomic.SwapInt32(&b.int32, 1) == 1
	} else {
		old = atomic.SwapInt32(&b.int32, 0) == 1
	}
	return
}
