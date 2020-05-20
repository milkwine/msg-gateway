package define

type Sender interface {
	Send(*Message) error
}

type Recorder interface {
	WrapMessage(m *Message) PostsendFunc
}

/*
1. sender application
2. http interface application
3. recorder application

http -> message -> wrap http response -> wrap recorder -> Send
*/

/*
Sender
1. real sender
2. abstract sender
  * queue
  * rate limit
  * switch sender

real sender -> rate limit sender -> retry sender -> queue sender -> switch sender

*/

/*
Recorder
1. db recorder
2. log recorder
*/
