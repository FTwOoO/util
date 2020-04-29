package timekit

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

var (
	timeFormat = "2006-01-02 15:04:05"
)

//获取某一天0点的时间
func GetStartTimeForDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

//获取某一天最后一时刻
func GetEndTimeForDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 1e9-1, time.Local)
}

//获取当前的系统时间戳（毫秒）
func CurrentTimestampInMils() uint64 {
	nowInMilis := time.Now().UnixNano() / 1e6
	return uint64(nowInMilis)
}

//获取当前的系统日期字符串2006-01-02格式
func CurrentDayString() string {
	dateString := time.Now().Local().Format("2006-01-02")
	return dateString
}

//获取当前的系统日期字符串2006-01-02 00:00:00格式
func CurrentDayTimeString() string {
	dateString := time.Now().Local().Format("2006-01-02 15:04:05")
	return dateString
}

func GetDayTimeString(t time.Time) string {
	dateString := t.Local().Format("2006-01-02 15:04:05")
	return dateString
}

func ParseDateTimeString(dayTimeStr string) (t time.Time, err error) {
	return time.ParseInLocation("2006-01-02 15:04:05", dayTimeStr, time.Local)
}

//获取昨天的系统日期字符串2006-01-02格式
func YesterdayDateString() string {
	nTime := time.Now()
	yesTime := nTime.AddDate(0, 0, -1)
	dateString := yesTime.Local().Format("2006-01-02")
	return dateString
}

//获取指定的日期字符串2006-01-02格式
func GetDateString(t time.Time) string {
	dateString := t.Local().Format("2006-01-02")
	return dateString
}

//获取最近7天字符串2006-01-02格式
func GetLast7DaysString() []string {
	var days []string
	n := 7
	now := time.Now()
	for i := 0; i < n; i++ {
		t := now.AddDate(0, 0, -i)
		dateString := t.Local().Format("2006-01-02")
		days = append(days, dateString)

	}
	return days
}

//获取指定的日期的月份字符串
func GetMonthString(t time.Time) string {
	dateString := t.Local().Format("2006-01")
	return dateString
}

//获取上一个月份字符串
func GetLastMonthString() string {
	return GetMonthString(time.Now().AddDate(0, -1, 0))
}

func ParseDateString(dayStr string) (t time.Time, err error) {
	return time.ParseInLocation("2006-01-02", dayStr, time.Local)
}

//时间是否在今天之后
func IsTimeAfterToday(t time.Time) bool {
	now := time.Now()
	if t.After(time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)) {
		return true
	}

	return false
}

func GetRandomDurationUpOrDown(duration time.Duration, deltaPercent float64) (ret time.Duration) {

	delta := time.Duration(rand.Int63n(int64(math.Round(float64(duration) * deltaPercent))))

	if rand.Intn(100)%2 == 1 {
		ret = duration + delta
	} else {
		ret = duration - delta
	}
	return
}

func GetRandomDuration(minDuration time.Duration, maxDuration time.Duration) time.Duration {
	randDuration := time.Duration(rand.Int63n(int64(maxDuration-minDuration))) + minDuration
	return randDuration
}

type TimeDurationSlice []time.Duration

func (p TimeDurationSlice) Len() int           { return len(p) }
func (p TimeDurationSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p TimeDurationSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p TimeDurationSlice) Sort()              { sort.Sort(p) }

//传入一个时长，在这个时长内随机时间内返回N个通知信号
func RandomTimeEvent(delay time.Duration, N int) <-chan time.Duration {
	ch := make(chan time.Duration)

	go func(ch chan time.Duration) {
		triggerPoint := make([]time.Duration, N)
		triggerPoint[0] = 0
		for i := 1; i < N-1; i += 1 {
			triggerPoint[i] = time.Duration(rand.Int63n(int64(delay)))
		}
		triggerPoint[N-1] = delay

		sort.Sort(TimeDurationSlice(triggerPoint))

		ch <- time.Duration(0)
		for i := 1; i < N; i += 1 {
			currentDelay := triggerPoint[i] - triggerPoint[i-1]
			if currentDelay < 0 {
				panic("???")
			}
			<-time.After(currentDelay)
			ch <- triggerPoint[i]
		}
		close(ch)
	}(ch)

	return ch
}
