package taskit

import (
	"testing"
	"time"
)

func TestRunOneByOne(t *testing.T) {
	RunOneByOne("testJob", func() (data interface{}, err error) {
		time.Sleep(1 * time.Second)
		t.Log("doing task ...")
		return
	},

		WithDelay(5*time.Second),
	)
}
