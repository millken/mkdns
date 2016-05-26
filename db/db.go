package db

import (
	"math/rand"
	"sync"

	"github.com/millken/mkdns/types"
)

type Db struct {
	sync.RWMutex
	m map[string]*types.Zone
}

const partition = 64

var (
	bucket [partition]Db
	seed   uint32
)

func Set(k string, z *types.Zone) {
	b := &bucket[hash(k)]
	b.Lock()
	b.m[k] = z
	b.Unlock()
}

func Get(k string) (z *types.Zone, ok bool) {
	b := &bucket[hash(k)]
	b.RLock()
	z, ok = b.m[k]
	b.RUnlock()
	return
}

func Delete(k string) {
	b := &bucket[hash(k)]
	if _, ok := b.m[k]; ok {
		delete(b.m, k)
	}
}

func hash(k string) uint32 {
	return djb33(k) % partition
}

// djb2 with better shuffling. 5x faster than FNV with the hash.Hash overhead.
func djb33(k string) uint32 {
	var (
		seed = rand.Uint32()
		l    = uint32(len(k))
		d    = 5381 + seed + l
		i    = uint32(0)
	)
	// Why is all this 5x faster than a for loop?
	if l >= 4 {
		for i < l-4 {
			d = (d * 33) ^ uint32(k[i])
			d = (d * 33) ^ uint32(k[i+1])
			d = (d * 33) ^ uint32(k[i+2])
			d = (d * 33) ^ uint32(k[i+3])
			i += 4
		}
	}
	switch l - i {
	case 1:
	case 2:
		d = (d * 33) ^ uint32(k[i])
	case 3:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
	case 4:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
		d = (d * 33) ^ uint32(k[i+2])
	}
	return d ^ (d >> 16)
}
