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
	// если получилось больше 1 - есть еще блокировки
	// ожидаем сингнал из канала (ch)
	if atomic.AddInt32(&mu.nlocks, 1) == 1 {
		return
	}

	// Ожидание сигнала из Unlock()
	<-mu.ch

}

func (mu *MuContest) LockChannel() <-chan struct{} {

	// if atomic.LoadInt32(&mu.nlocks) == 0 {
	// 	mu.ch <- mu.sig
	// }

	if atomic.AddInt32(&mu.nlocks, 1) == 1 {
		mu.ch <- mu.sig
	}

	return mu.ch

}

func (mu *MuContest) Unlock() {

	// Уменьшение счетчика блокировок
	// если есть еще блокировки отправляем сигнал в канал (ch).
	if atomic.AddInt32(&mu.nlocks, -1) == 0 {
		return
	}

	if atomic.LoadInt32(&mu.nlocks) < 0 {
		panic("unlocked")
	}

	// Сигнал о разблокировке
	mu.ch <- mu.sig

}

func New() Mutex {
	return &MuContest{ch: make(chan struct{}, 1)}
}
