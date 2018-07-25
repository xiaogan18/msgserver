package persistence
import(
	"time"
	"github.com/xiaogan18/msgserver/queue"
	"github.com/xiaogan18/msgserver/serialize"
	"errors"
)
var(
	NotFoundMsg=errors.New("msg is not found")
)
type MsgContainer interface{
	Get(id string) (*OfflineMsg,error)
	GetUserMsg(userID string) ([]*OfflineMsg,error)
	Put(*OfflineMsg)
}
type OfflineMsg struct{
	queue.Message
	KeepLiveTime time.Time
}

func CreateMsgContainer(t string,params ...interface{}) MsgContainer{
	switch(t){
	case "redis":
		p:=params[0].(*RedisOptions)
		c:=RedisContainer{
			serializer:serialize.CreateSerializer("json"),
			option:p,
		}
		c.try()
		return &c
	default:
		c:=MemoryContainer{
			msgMap:make(map[string]*OfflineMsg,0),
			userMsgMap:make(map[string][]string,0),
		}
		c.gc(1000*120)
		return &c
	}
}