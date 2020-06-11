package mem

import (
	"github.com/golang/protobuf/proto"
	"reflect"
	"testing"
)

type A struct {
	name   string
	avatar string
}

func TestMessagePool_PutObject(t *testing.T) {
	PoolPutObject(&A{name: "1", avatar: "https://baidu.com"})
	v := PoolGetObject((*A)(nil))
	t.Logf("%v %v", reflect.TypeOf(v), v)
}


func TestReset(t *testing.T) {
	buffer1 := PoolGetObject((*proto.Buffer)(nil)).(*proto.Buffer)
	err := buffer1.EncodeFixed32(11)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", buffer1.Bytes())

	PoolPutObject(buffer1)
	buffer2 := PoolGetObject((*proto.Buffer)(nil)).(*proto.Buffer)

	if len(buffer2.Bytes()) != 0 {
		t.Fatalf("Reset() call of *proto.Buffer fail")
	}
}


func TestResetForProtoMessage(t *testing.T) {
	msg := PoolGetObject((*TestMsg)(nil)).(*TestMsg)
	msg.RoomId = "10004"
	PoolPutObject(msg)
	msg = PoolGetObject((*TestMsg)(nil)).(*TestMsg)
	t.Log(msg)
	if msg.RoomId != "" {
		t.Fatalf("Reset() call of proto.Message fail")
	}
}