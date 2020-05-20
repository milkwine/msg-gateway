package define

import (
	"msg-gateway/async"
)

type PostsendFunc func(e error)

type Message struct {
	ID        string
	Type      MsgType
	HandleStr string
	Content   string
	postFn    []PostsendFunc
}

func (m *Message) AddPostFunc(fn PostsendFunc) {
	if m.postFn == nil {
		m.postFn = make([]PostsendFunc, 0)
	}
	m.postFn = append(m.postFn, fn)
}

func (m *Message) InvokePostFn(e error) {

	for _, fn := range m.postFn {
		f := fn
		async.Async(func() { f(e) })
	}
}

type MsgType uint

const (
	TSMS MsgType = iota
)

func (t MsgType) String() string {
	switch t {
	case TSMS:
		return "SMS"
	default:
		return "UNKNOWN"
	}
}
