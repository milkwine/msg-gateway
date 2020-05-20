package main

import (
	"msg-gateway/async"
	"msg-gateway/define"
	"msg-gateway/recorder"
	"msg-gateway/sender"
)

func main() {
	send := sender.NewDebugSender()
	record := &recorder.Debug{}

	go send.Start()

	m := define.Message{
		ID:        "testID",
		Type:      1,
		HandleStr: "123123",
		Content:   "test content",
	}

	m.AddPostFunc(record.WrapMessage(&m))
	m.AddPostFunc(func(e error) {
		panic("test panic!")
	})

	send.Send(&m)
	send.Stop()
	async.Pool.Stop()
}
