package taskit

import (
	"github.com/FTwOoO/util/logging"
	"time"
)

type runOption func(cf *supervisorConfig)

func WithLock(aquireLock func() (success bool, err error)) runOption {
	return func(cf *supervisorConfig) {
		cf.aquireLock = aquireLock
	}
}

func WithLockRelease(lockRelease func() (err error)) runOption {
	return func(cf *supervisorConfig) {
		cf.lockRelease = lockRelease
	}
}

func WithDelay(duration time.Duration) runOption {
	return func(cf *supervisorConfig) {
		cf.delay = func() time.Duration {
			return duration
		}
	}
}

type supervisorConfig struct {
	aquireLock  func() (success bool, err error)
	lockRelease func() (err error)
	delay       func() time.Duration
}

func RunOneByOne(
	jobName string,
	f func() (data interface{}, err error),
	opts ...runOption) {

	cf := &supervisorConfig{}
	for _, optFunc := range opts {
		optFunc(cf)
	}

	for {
		if cf.aquireLock != nil {
			success, err := cf.aquireLock()

			if !success || err != nil {
				if err != nil {
					logging.Log.Errorw(logging.KeyEvent, "fetchLockFail", "err", err)
				}

				if cf.delay != nil {
					<-time.After(cf.delay() / 10)
				} else {
					<-time.After(10 * time.Second)
				}
				continue
			}
		}

		defaultTaskManager.Run(jobName, f)
		_, err := defaultTaskManager.Wait(jobName)
		if err != nil {
			logging.Log.Infow(logging.KeyEvent, "jobFail", "job", jobName, "err", err)
			if cf.lockRelease != nil {
				logging.Log.Infow(logging.KeyEvent, "jobLockRelease", "job", jobName)
				_ = cf.lockRelease()
			}
		}

		if cf.delay != nil {
			logging.Log.Infow(logging.KeyEvent, "sleepForJob", "job", jobName, "sleepDuration", cf.delay())
			<-time.After(cf.delay())
		}
	}
}
