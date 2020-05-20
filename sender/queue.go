package sender

import (
	"errors"
	"msg-gateway/define"
	"sync"
	"time"
)

type QueueSender struct {
	s       define.Sender
	queue   chan *define.Message
	running bool
	sig     chan *struct{}
	lock    sync.RWMutex
	timeout time.Duration
}

//s should start
func NewQueueSender(s define.Sender, qLen int) *QueueSender {

	return &QueueSender{
		s:       s,
		queue:   make(chan *define.Message, qLen),
		sig:     make(chan *struct{}),
		running: false,
		lock:    sync.RWMutex{},
		timeout: time.Millisecond * 100,
	}
}

func (q *QueueSender) Start() error {

	{
		q.lock.Lock()

		if q.running {

			q.lock.Unlock()
			return errors.New("already start")
		}
		q.running = true

		q.lock.Unlock()
	}

	for {

		m := <-q.queue
		if m == nil {
			break
		}

		q.s.Send(m)
	}

	q.sig <- nil

	return nil
}

func (q *QueueSender) Send(m *define.Message) error {

	q.lock.RLocker().Lock()
	defer q.lock.RLocker().Unlock()

	if !q.running {
		return errors.New("queueSender is not running")
	}

	ticker := time.NewTicker(q.timeout)
	defer ticker.Stop()

	select {
	case q.queue <- m:
		return nil
	case <-ticker.C:
		return errors.New("enqueue timeout!")
	}
}

func (q *QueueSender) Stop() error {

	q.lock.Lock()
	defer q.lock.Unlock()

	if q.running == false {
		return nil
	}

	close(q.queue)

	<-q.sig

	q.running = false
	return nil
}
