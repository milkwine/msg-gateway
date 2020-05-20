package recorder

import (
	"fmt"
	"msg-gateway/define"
)

type Debug struct {
}

func (_ *Debug) WrapMessage(m *define.Message) define.PostsendFunc {

	fmt.Printf("[Recorder Log] Msg %s Begin send.\n", m.ID)

	return func(e error) {
		if e == nil {
			fmt.Printf("[Recorder Log] Msg %s send suc.\n", m.ID)

		} else {
			fmt.Printf("[Recorder Log] Msg %s send fail(%s).\n", m.ID, e.Error())
		}

	}
}
