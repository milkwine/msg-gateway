package sender

import (
	"fmt"
	"msg-gateway/define"
	"msg-gateway/recorder"
	"testing"
	"time"

	"github.com/k0kubun/pp"
)

func TestQueue(t *testing.T) {
	debug := NewDebugSender()
	q := NewQueueSender(debug, 50)
	go q.Start()

	rq := NewRetrySender(q, 5, time.Millisecond*300)
	pp.Println(rq)

	record := &recorder.Debug{}

	for i := 0; i < 100; i++ {

		m := define.Message{
			ID:        fmt.Sprintf("testID-%d", i),
			Type:      1,
			HandleStr: "123123",
			Content:   "test content",
		}
		m.AddPostFunc(record.WrapMessage(&m))
		err := rq.Send(&m)
		if err != nil {

			pp.Println(err)

		}
	}

	q.Stop()
}
