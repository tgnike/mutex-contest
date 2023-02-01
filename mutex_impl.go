package contest

import (
	"sync/atomic"
)

type MuContest struct {
	nlocks int32
	ch     chan struct{}
	sig    struct{}
}

func (mu *MuContest) Lock() {

	// Увеличение количества блокировок
	nlocks := atomic.AddInt32(&mu.nlocks, 1)

	// Если больше одной блокировки ждем...
	if nlocks != 1 {
		<-mu.ch
	}

}

func (mu *MuContest) LockChannel() <-chan struct{} {

	ch := make(chan struct{}, 1)

	// Если ничего не заблокировано, отправляем сигнал в канал.
	if mu.nlocks == 0 {
		ch <- mu.sig
	}

	return ch

}

func (mu *MuContest) Unlock() {

	if mu.nlocks == 0 {
		return
	}

	nLocks := atomic.AddInt32(&mu.nlocks, -1)

	// Если не все блокировки сняты, отправляем сингнал о разблокировке
	if nLocks != 0 {
		mu.ch <- mu.sig
	}

}

func New() Mutex {
	return &MuContest{ch: make(chan struct{}, 1)}
}
