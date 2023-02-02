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
	if atomic.AddInt32(&mu.nlocks, 1) == 1 {
		return
	}

	<-mu.ch

}

func (mu *MuContest) LockChannel() <-chan struct{} {

	//ch := make(chan struct{}, 1)

	if atomic.LoadInt32(&mu.nlocks) == 0 {
		mu.ch <- mu.sig
	}

	// Если ничего не заблокировано, отправляем сигнал в канал.
	// if mu.nlocks == 0 {
	// 	mu.ch <- mu.sig
	// }

	return mu.ch

}

func (mu *MuContest) Unlock() {

	// if mu.nlocks == 0 {
	// 	return
	// }

	if atomic.AddInt32(&mu.nlocks, -1) == 0 {
		return
	}

	// Если не все блокировки сняты, отправляем сингнал о разблокировке
	mu.ch <- mu.sig

}

func New() Mutex {
	return &MuContest{ch: make(chan struct{}, 1)}
}
