package rutine

import (
	"github.com/panjf2000/ants"
	"sync"
)

var pool *ants.Pool

var once sync.Once

func GetPoolInstance() *ants.Pool {
	once.Do(func() {
		p, err := ants.NewPool(10)
		if err != nil {
			panic(err)
		}
		pool = p
	})
	return pool
}

func Submit(fun func()) {
	GetPoolInstance().Submit(fun)
}
