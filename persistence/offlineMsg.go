package persistence
import(
	"time"
	"msgserver/queue"
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