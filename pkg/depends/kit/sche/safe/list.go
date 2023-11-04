package safe

import (
	"container/list"
	"sync"
)

type Element = list.Element

type List struct {
	mu *sync.RWMutex
	list.List
}

func NewList(v ...interface{}) *List {
	ret := (&List{mu: &sync.RWMutex{}}).Init()
	for i := range v {
		ret.List.PushBack(v[i])
	}
	return ret
}

func (l *List) lock()    { l.mu.Lock() }
func (l *List) unlock()  { l.mu.Unlock() }
func (l *List) rLock()   { l.mu.RLock() }
func (l *List) rUnlock() { l.mu.RUnlock() }

func (l *List) Init() *List {
	l.List.Init()
	return l
}

func (l *List) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.List.Len()
}

func (l *List) PushFront(v interface{}) *Element {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.List.PushFront(v)
}

func (l *List) PushFrontN(v ...interface{}) {
	l.lock()
	defer l.unlock()
	for i := range v {
		l.List.PushFront(v[i])
	}
}

func (l *List) PushFrontList(lst *List) {
	l.lock()
	defer l.unlock()
	l.List.PushFrontList(&lst.List)
}

func (l *List) PushBack(v interface{}) *Element {
	l.lock()
	defer l.unlock()
	return l.List.PushBack(v)
}

func (l *List) PushBackN(v ...interface{}) {
	l.lock()
	defer l.unlock()
	for i := range v {
		l.List.PushBack(v[i])
	}
}

func (l *List) PushBackList(lst *List) {
	l.lock()
	defer l.unlock()
	l.List.PushBackList(&lst.List)
}

func (l *List) PopFront() interface{} {
	l.lock()
	defer l.unlock()
	if e := l.List.Front(); e != nil {
		return l.List.Remove(e)
	}
	return nil
}

func (l *List) PopFrontN(n int) (ret []interface{}) {
	l.lock()
	defer l.unlock()
	for i := 0; i < n; i++ {
		if e := l.List.Front(); e != nil {
			ret = append(ret, l.List.Remove(e))
			continue
		}
		break
	}
	return
}

func (l *List) FrontN(n int) (ret []interface{}) {
	l.lock()
	defer l.unlock()
	iter := l.List.Front()
	for i := 0; i < n; i++ {
		if iter == nil {
			break
		}
		ret = append(ret, iter.Value)
		iter = iter.Next()
	}
	return ret
}

func (l *List) PopBack() interface{} {
	l.lock()
	defer l.unlock()
	if e := l.List.Back(); e != nil {
		return l.List.Remove(e)
	}
	return nil
}

func (l *List) PopBackN(n int) (ret []interface{}) {
	l.lock()
	defer l.unlock()
	for i := 0; i < n; i++ {
		if e := l.List.Back(); e != nil {
			ret = append(ret, l.List.Remove(e))
			continue
		}
		break
	}
	return
}

func (l *List) Front() *Element {
	l.rLock()
	defer l.rUnlock()
	return l.List.Front()
}

func (l *List) Back() *Element {
	l.rLock()
	defer l.rUnlock()
	return l.List.Back()
}

func (l *List) FrontValue() interface{} {
	l.rLock()
	defer l.rUnlock()
	if e := l.List.Front(); e != nil {
		return e.Value
	}
	return nil
}

func (l *List) BackValue() interface{} {
	l.rLock()
	defer l.rUnlock()
	if e := l.List.Back(); e != nil {
		return e.Value
	}
	return nil
}

func (l *List) MoveToFront(e *Element) {
	l.lock()
	defer l.unlock()
	l.List.MoveToFront(e)
}

func (l *List) MoveToBack(e *Element) {
	l.lock()
	defer l.unlock()
	l.List.MoveToBack(e)
}

func (l *List) MoveBefore(e, mark *Element) {
	l.lock()
	defer l.unlock()
	l.List.MoveBefore(e, mark)
}

func (l *List) MoveAfter(e, mark *Element) {
	l.lock()
	defer l.unlock()
	l.List.MoveAfter(e, mark)
}

func (l *List) InsertBefore(v interface{}, e *Element) *Element {
	l.lock()
	defer l.unlock()
	return l.List.InsertBefore(v, e)
}

func (l *List) InsertAfter(v interface{}, e *Element) *Element {
	l.lock()
	defer l.unlock()
	return l.List.InsertAfter(v, e)
}

func (l *List) Remove(e ...*Element) {
	l.lock()
	defer l.unlock()
	for i := range e {
		l.List.Remove(e[i])
	}
}

func (l *List) Clear() {
	l.lock()
	defer l.unlock()
	for l.Len() > 0 {
		l.List.Remove(l.List.Front())
	}
}

func (l *List) Elements() (ret []interface{}) {
	l.Range(func(e *Element) bool {
		ret = append(ret, e.Value)
		return true
	})
	return ret
}

func (l *List) RElements() (ret []interface{}) {
	l.ReverseRange(func(e *Element) bool {
		ret = append(ret, e.Value)
		return true
	})
	return ret
}

func (l *List) Range(f func(*Element) bool) {
	l.rLock()
	defer l.rUnlock()
	for e := l.List.Front(); e != nil; {
		if f(e) {
			e = e.Next()
			continue
		}
		break
	}
}

func (l *List) ReverseRange(f func(*Element) bool) {
	l.rLock()
	defer l.rUnlock()
	for e := l.List.Back(); e != nil; {
		if f(e) {
			e = e.Prev()
			continue
		}
		break
	}
}

func (l *List) RLockDo(f func(*list.List)) {
	l.rLock()
	defer l.rUnlock()
	f(&l.List)
}

func (l *List) LockDo(f func(*list.List)) {
	l.lock()
	defer l.unlock()
	f(&l.List)
}
