package contest

import (
	"sync/atomic"
)

const lock int32 = 1

type MuContest struct {
	nlocks int32
	ch     chan struct{}
	sig    struct{}
}

// Lock Блокирует мьютекс
// Если блокировка уже установлена - блокирует горутину до разблокировки
func (mu *MuContest) Lock() {

	// Увеличение счетчика блокировок
	// если получилось больше 1 - есть еще блокировки
	// ожидаем сингнал из канала (ch)
	if atomic.AddInt32(&mu.nlocks, lock) == 1 {
		return
	}

	// Ожидание сигнала из канала
	<-mu.ch

}

// LockChannel блокировка канала
func (mu *MuContest) LockChannel() <-chan struct{} {

	if atomic.AddInt32(&mu.nlocks, lock) == 1 {
		mu.ch <- mu.sig
	}

	return mu.ch

}

// Unlock разблокирует мьютекс
func (mu *MuContest) Unlock() {

	// Уменьшение счетчика блокировок
	// если есть еще блокировки отправляем сигнал в канал (ch).
	if atomic.AddInt32(&mu.nlocks, -lock) == 0 {
		return
	}

	// Паника: разблокировка без блокировки
	if atomic.LoadInt32(&mu.nlocks) < 0 {
		panic("attempt to unlock when there are no locks")
	}

	// Сигнал о разблокировке
	mu.ch <- mu.sig

}

func New() Mutex {
	return &MuContest{ch: make(chan struct{}, 1)}
}
