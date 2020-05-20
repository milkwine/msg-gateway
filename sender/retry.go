package sender

import (
	"msg-gateway/define"
	"time"
)

type RetrySender struct {
	s        define.Sender
	closed   chan struct{}
	maxErr   int
	interval time.Duration
}

//s should start
func NewRetrySender(s define.Sender, maxErr int, interval time.Duration) *RetrySender {

	return &RetrySender{
		s:        s,
		closed:   make(chan struct{}),
		maxErr:   maxErr,
		interval: interval,
	}
}

func (q *RetrySender) Send(m *define.Message) error {

	run := 0
	for {
		run += 1

		err := q.s.Send(m)

		if err == nil {
			return nil
		}

		if run >= q.maxErr {
			return err
		}

		time.Sleep(q.interval)
	}
}
