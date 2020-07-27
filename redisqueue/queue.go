package redisqueue

import (
	"fmt"
	"github.com/go-redis/redis"
	"gitlab.livedev.shika2019.com/go/common/cfg"
	"gitlab.livedev.shika2019.com/go/common/rediskit"
	"strconv"
	"time"
)

// Queue holds a reference to a redis connection and a queue name.
type Queue struct {
	cli  *rediskit.RedisClient
	Name string
}

// New defines a new Queue
func New(queueName string, cf *cfg.RedisConfig) *Queue {
	return &Queue{
		cli:  rediskit.NewRedisClient(cf),
		Name: queueName,
	}
}

// Push pushes a single job on to the queue. The job string can be any format, as the queue doesn't really care.
func (q *Queue) Push(job string) (bool, error) {
	return q.Schedule(job, time.Now())
}

// Schedule schedule a job at some point in the future, or some point in the past. Scheduling a job far in the past is the same as giving it a high priority, as jobs are popped in order of due date.
func (q *Queue) Schedule(job string, when time.Time) (bool, error) {
	score := when.UnixNano() / 1e6
	res, err := q.cli.ZAdd(q.Name, redis.Z{Member: job, Score: float64(score)}).Result()
	if err != nil {
		return false, err
	}
	// _, err := addTaskScript.Do(q.c, job)
	return res == 1, err

}

// Pending returns the count of jobs pending, including scheduled jobs that are not due yet.
func (q *Queue) Pending() (int64, error) {
	return q.cli.ZCard(q.Name).Result()
}

// FlushQueue removes everything from the queue. Useful for testing.
func (q *Queue) FlushQueue() error {
	return q.cli.Del(q.Name).Err()
}

// Pop removes and returns a single job from the queue. Safe for concurrent use (multiple goroutines must use their own Queue objects and redis connections)
func (q *Queue) Pop() (string, error) {
	jobs, err := q.PopJobs(1)
	if err != nil {
		return "", err
	}
	if len(jobs) == 0 {
		return "", nil
	}
	return jobs[0], nil
}

// PopJobs returns multiple jobs from the queue. Safe for concurrent use (multiple goroutines must use their own Queue objects and redis connections)
func (q *Queue) PopJobs(limit int) ([]string, error) {
	IncrByXX := redis.NewScript(`
		local name = KEYS[1]
		local timestamp = KEYS[2]
    local limit = KEYS[3]
		local results = redis.call('zrangebyscore', name, '-inf', timestamp, 'LIMIT', 0, limit)
		if table.getn(results) > 0 then
			redis.call('zrem', name, unpack(results))
		end
    return results
  `)

	ret, err := IncrByXX.Run(q.cli, []string{q.Name, fmt.Sprintf("%d", time.Now().UnixNano()/1e6), strconv.Itoa(limit)}).Result()
	if err != nil {
		return nil, err
	}

	vals := []string{}
	for _, obj := range ret.([]interface{}) {
		if _, ok := obj.(string); ok {
			vals = append(vals, obj.(string))
		}
	}

	return vals, nil
}
