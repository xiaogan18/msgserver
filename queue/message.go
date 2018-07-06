package queue


type Message struct{
	MsgId string
	MsgType MessageType
	Content interface{}
	To string
}
//消息类型
type MessageType int
const(
	// 针对消息 发送给单个终端
	Message_ToOne MessageType=iota
	// 广播消息
	Message_ToAll
)
