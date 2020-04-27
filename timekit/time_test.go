package timekit

import (
	"testing"
	"time"
)

func Test_GetStartTimeForDay(t *testing.T) {
	now := time.Now()
	t.Logf("start of today:%s", GetStartTimeForDay(now).String())
	t.Logf("start of today:%s", GetEndTimeForDay(now).String())

	diff := GetEndTimeForDay(now).Sub(GetStartTimeForDay(now))

	if (diff + time.Duration(1*time.Nanosecond)) != time.Hour*24 {
		t.Fatalf("end time - start time = %v", diff)
	}
}

func TestGetRandomDuration(t *testing.T) {
	min := 1 * time.Second
	max := 10 * time.Second
	for i := 0; i < 100; i++ {
		out := GetRandomDuration(min, max)

		if out < min || out > max {
			t.FailNow()
		}
		t.Log(out)

	}
}

func TestRandomTimeEvent(t *testing.T) {

	n := 15
	ch := RandomTimeEvent(1*time.Second, n)
	for x := range ch {
		t.Log(x)
	}
}
