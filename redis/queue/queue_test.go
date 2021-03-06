package queue

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/laonsx/gamelib/codec"
	"github.com/laonsx/gamelib/gofunc"
	"github.com/laonsx/gamelib/redis"
	"github.com/laonsx/gamelib/task"
)

func init() {

	uinfoTask = task.New(32, syncUserInfoTask)

	redis.InitRedis(codec.MsgPack, codec.UnMsgPack, redis.NewRedisConf("queue", "127.0.0.1", "6379", 0))

	RegisterQueueHandler("SyncUserInfo", syncUserInfoHandle)
}

var uinfoTask *task.Task

type QMsg struct {
	Type   string
	Uid    uint64
	Params interface{}
}

func syncUserInfoTask(v interface{}) {
	defer gofunc.PrintPanic()

	data := v.([]byte)
	qmsg := new(QMsg)
	err := json.Unmarshal(data, qmsg)
	if err != nil {
		fmt.Println("SyncUserInfo handle", err.Error())
		return
	}

	fmt.Println("queue handle", qmsg.Uid)

	switch qmsg.Type {
	case "Battle":
	case "Name":

		fmt.Println("sync user name:", qmsg)
	case "LoginTime":
	}
}

func syncUserInfo(uid uint64, atype string, params ...interface{}) {
	qmsg := new(QMsg)
	qmsg.Type = atype
	qmsg.Uid = uid
	if len(params) > 0 {
		qmsg.Params = params[0]
	}

	data, err := json.Marshal(qmsg)
	if err != nil {
		fmt.Println("SyncUserInfo uid=", uid, "err=", err.Error())
		return
	}

	err = QPush("SyncUserInfo", data)
	if err != nil {
		fmt.Println("SyncUserInfo uid=", uid, "err=", err.Error())
	}
}

func syncUserInfoHandle(data []byte) {
	err := uinfoTask.SendMsg(data)
	if err != nil {
		fmt.Println("syncUserInfoHandle", err.Error())
	}
}

func TestQueue(t *testing.T) {

	for i := 100; i < 200; i++ {

		syncUserInfo(uint64(i), "Name", fmt.Sprintf("%d_name_%d", i, i))
	}
}
