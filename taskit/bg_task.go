package taskit

import (
	"gitlab.livedev.shika2019.com/go/util/logging"
	"runtime/debug"
)

func RunInBackend(f func()) {
	defer func() {
		if x := recover(); x != nil {
			logging.Log.Fatalw(logging.KeyEvent, "panic", "msg", x, "stack", string(debug.Stack()))
		}
	}()
	f()
}
