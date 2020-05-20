package async

import (
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var Pool *AsyncPool

func init() {
	Pool = &AsyncPool{}
	Pool.StartWithSize(10, 100)
}

func Async(fn PoolFunc) {

	f := wrapFunc(fn)
	err := Pool.Enqueue(f)

	if err != nil {
		// backoff to sync fuc
		log.Printf(err.Error())
		f()
	}
}

func wrapFunc(fn PoolFunc) PoolFunc {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in poolFunc", r)
			}
		}()
		fn()
	}
}

type PoolFunc func()

type AsyncPool struct {
	workers int
	queue   chan PoolFunc
	start   bool
	once    sync.Once
	wg      *sync.WaitGroup
}

func (p *AsyncPool) Start(workerNum int) error {

	return p.StartWithSize(workerNum, workerNum*2)
}

func (p *AsyncPool) StartWithSize(workerNum int, queSize int) error {

	if p.start {
		return errors.New("Has start already")
	}

	p.wg = &sync.WaitGroup{}
	p.workers = workerNum
	p.queue = make(chan PoolFunc, queSize)

	for i := 0; i < workerNum; i++ {

		p.wg.Add(1)
		go func(id int, queue chan PoolFunc, wg *sync.WaitGroup) {
			defer wg.Done()

			log.Printf("[AsyncPool] worker %d start.\n", id)

			for {
				fn := <-queue
				if fn == nil {
					break
				}
				fn()
			}
			log.Printf("[AsyncPool] worker %d exit.\n", id)
		}(i+1, p.queue, p.wg)

	}
	p.start = true
	return nil
}

func (p *AsyncPool) Enqueue(fn PoolFunc) error {

	if !p.start {
		return errors.Errorf("pool not start")
	}

	tick := time.NewTicker(time.Microsecond * 100)

	select {
	case p.queue <- fn:
		return nil
	case <-tick.C:
		return errors.Errorf("queue is full, len: %d", len(p.queue))
	}
}

func (p *AsyncPool) Stop() {

	p.once.Do(func() {
		close(p.queue)
		//for i := 0; i < p.workers; i++ {
		//	p.queue <- nil
		//}
	})

	p.start = false
	p.wg.Wait()
}
