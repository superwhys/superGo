package goroutinePool

import (
	"errors"
	"github.com/superwhys/superGo/superLog"
	"sync"
)

type Pool struct {
	chanQueue chan struct{}
	wg        *sync.WaitGroup
}

func NewPool(size int) *Pool {
	if size <= 0 {
		size = 1
	}
	return &Pool{
		chanQueue: make(chan struct{}, size),
		wg:        &sync.WaitGroup{},
	}
}

func (p *Pool) Add(delta int) {
	if delta < 0 {
		superLog.PanicError(errors.New("delta is less than 0"))
	}
	for i := 0; i < delta; i++ {
		p.chanQueue <- struct{}{}
	}
	p.wg.Add(delta)
}

func (p *Pool) Done() {
	<-p.chanQueue
	p.wg.Done()
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
