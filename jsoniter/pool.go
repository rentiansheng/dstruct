package jsoniter

import (
	"sync"
)

// IteratorPool a thread safe pool of iterators with same configuration
type IteratorPool interface {
	BorrowIterator(data []byte) *Iterator
	ReturnIterator(iter *Iterator)
}

var iteratorPool = sync.Pool{
	New: func() interface{} {
		return NewIterator()
	},
}

func BorrowIterator(data []byte) *Iterator {
	iter := iteratorPool.Get().(*Iterator)
	iter.ResetBytes(data)
	return iter
}

func ReturnIterator(iter *Iterator) {
	iter.Error = nil
	iter.Attachment = nil
	iteratorPool.Put(iter)
}
