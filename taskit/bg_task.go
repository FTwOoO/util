package taskit

import (
	"github.com/rexue2019/util/logging"
	"runtime/debug"
)

func RunInBackend(f func()) {
	go func() {
		defer func() {
			if x := recover(); x != nil {
				logging.Log.Fatalw(logging.KeyEvent, "panic", "msg", x, "stack", string(debug.Stack()))
			}
		}()
		f()
	}()
}
