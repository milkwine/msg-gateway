package sender

import (
	"fmt"
	"msg-gateway/define"
	"sync"
	"time"
)

type DebugSender struct {
	stop   chan struct{}
	closed chan struct{}
	once   sync.Once
}

func NewDebugSender() *DebugSender {

	return &DebugSender{
		stop:   make(chan struct{}),
		closed: make(chan struct{}),
	}
}

func (s *DebugSender) Start() error {

	defer close(s.closed)

	<-s.stop

	return nil
}

func (s *DebugSender) Stop() error {
	s.once.Do(func() { close(s.stop) })
	<-s.closed
	return nil
}

func (s *DebugSender) Send(msg *define.Message) error {

	fmt.Printf("[debug sender] ID: %s, Type: %d, HandleStr: %s, Content: %s\n", msg.ID, msg.Type, msg.HandleStr, msg.Content)
	time.Sleep(time.Second * 1)

	msg.InvokePostFn(nil)

	return nil
}
