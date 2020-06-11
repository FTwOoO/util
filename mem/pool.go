package mem

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"gitlab.livedev.shika2019.com/go/common/logging"
	"reflect"
	"sync"
)

var _messagePool *messagePool

func init() {
	_messagePool = newMessagePool()
}

type messagePool struct {
	poolsLock *sync.RWMutex
	pools     map[reflect.Type]*sync.Pool
}

func newMessagePool() *messagePool {
	return &messagePool{
		poolsLock: &sync.RWMutex{},
		pools:     map[reflect.Type]*sync.Pool{},
	}
}

func (this *messagePool) GetByObject(msg interface{}) (ret interface{}) {
	t := reflect.TypeOf(msg)
	if t == nil {
		logging.Log.Errorw(logging.KeyScope, "memPool", logging.KeyMsg, "TypeOf nil")
		return
	}

	return this.GetByType(t)
}

func (this *messagePool) createPool(t reflect.Type) {
	this.poolsLock.Lock()
	pool := &sync.Pool{
		New: func() interface{} {
			return reflect.New(t).Interface()
		},
	}
	this.pools[t] = pool
	this.poolsLock.Unlock()
}

func (this *messagePool) GetByType(t reflect.Type) (ret interface{}) {

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	this.poolsLock.RLock()
	pool, ok := this.pools[t]
	this.poolsLock.RUnlock()

	if !ok {
		this.createPool(t)
		this.poolsLock.RLock()
		pool = this.pools[t]
		this.poolsLock.RUnlock()
	}

	ret = pool.Get()
	if ret == nil {
		logging.Log.Fatalw(logging.KeyScope, "memPool", logging.KeyMsg, "invalid type", "type", t)
	}
	if v, ok := ret.(proto.Message); ok {
		v.Reset()
	} else if v, ok := ret.(*bytes.Buffer); ok {
		v.Reset()
	} else if v, ok := ret.(*proto.Buffer); ok {
		v.Reset()
	}
	return
}

func (this *messagePool) PutObject(msg interface{}) {
	t := reflect.TypeOf(msg)
	if t == nil {
		logging.Log.Errorw(logging.KeyScope, "memPool", logging.KeyMsg, "TypeOf nil")
		return
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	this.poolsLock.RLock()
	pool, ok := this.pools[t]
	this.poolsLock.RUnlock()

	if !ok {
		logging.Log.Infow(
			logging.KeyScope, "memPool",
			logging.KeyEvent, "poolNotExist",
			"type", t)

		this.createPool(t)
		this.poolsLock.RLock()
		pool = this.pools[t]
		this.poolsLock.RUnlock()
	}
	pool.Put(msg)
}

func PoolGetObject(msg interface{}) (interface{}) {
	return _messagePool.GetByObject(msg)
}

func PoolGetObjectByType(t reflect.Type) (interface{}) {
	return _messagePool.GetByType(t)
}

func PoolPutObject(msg interface{}) {
	_messagePool.PutObject(msg)
}
