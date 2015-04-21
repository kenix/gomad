package worker

import (
	"sync"
)

type Worker interface {
	Task(int)
}

type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

func New(maxGoroutines int) *Pool {
	p := Pool{work: make(chan Worker)}
	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func(wid int) {
			for w := range p.work {
				w.Task(wid)
			}
			p.wg.Done()
		}(i)
	}

	return &p
}

func (p *Pool) Run(w Worker) {
	p.work <- w
}

func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
