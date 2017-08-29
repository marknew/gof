/**
 * Copyright 2017 by mark zhang zhangjinghui@dafy.com
 * name :data synchronized server from beijing dafy
 * author : mark zhang
 * date : 2017-02-11 12:28
 * description :临时对象池
 * history :
 */
package util

import "sync"

var DEFAULT_SYNC_POOL *SyncPool

func NewPool() *SyncPool {
	DEFAULT_SYNC_POOL = NewSyncPool(
		5,
		30000,
		2,
	)
	return DEFAULT_SYNC_POOL
}

func Alloc(size int) []interface{} {
	return DEFAULT_SYNC_POOL.Alloc(size)
}

func Free(mem []interface{}) {
	DEFAULT_SYNC_POOL.Free(mem)
}

// SyncPool is a sync.Pool base slab allocation memory pool
type SyncPool struct {
	classes     []sync.Pool
	classesSize []int
	minSize     int
	maxSize     int
}

func NewSyncPool(minSize, maxSize, factor int) *SyncPool {
	n := 0
	for chunkSize := minSize; chunkSize <= maxSize; chunkSize *= factor {
		n++
	}
	pool := &SyncPool{
		make([]sync.Pool, n),
		make([]int, n),
		minSize, maxSize,
	}
	n = 0
	for chunkSize := minSize; chunkSize <= maxSize; chunkSize *= factor {
		pool.classesSize[n] = chunkSize
		pool.classes[n].New = func(size int) func() interface{} {
			return func() interface{} {
				buf := make([]interface{}, size)
				return &buf
			}
		}(chunkSize)
		n++
	}
	return pool
}

func (pool *SyncPool) Alloc(size int) []interface{} {
	if size <= pool.maxSize {
		for i := 0; i < len(pool.classesSize); i++ {
			if pool.classesSize[i] >= size {
				mem := pool.classes[i].Get().(*[]interface{})
				// return (*mem)[:size]
				return (*mem)[:0]
			}
		}
	}
	return make([]interface{}, 0, size)
}

func (pool *SyncPool) Free(mem []interface{}) {
	if size := cap(mem); size <= pool.maxSize {
		for i := 0; i < len(pool.classesSize); i++ {
			if pool.classesSize[i] >= size {
				pool.classes[i].Put(&mem)
				return
			}
		}
	}
}
