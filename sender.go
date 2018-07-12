package msgserver
import(
	"msgserver/queue"
	"msgserver/pool"
	"msgserver/persistence"
	"time"
	"fmt"
)

type Sender struct{
	pool pool.Pool
	queue queue.Queue
	// tcp代理类
	tcpProxy *TcpProxy
	// 消息持久化容器
	Container persistence.MsgContainer
	// 消息发送失败时回调
	FailedCallback func(error)  
	// 失败重试次数
	ResendTimes int 
	// 重试间隔
	ResendInterval int 
}
// 发送消息
func(this *Sender) SendMessage(msg interface{},user string) error{
	m:=queue.Message{
		MsgId:fmt.Sprintf("1%s",Guid()),
		MsgType:queue.Message_ToOne,
		Content:msg,
		To:user,
	}
	return this.queue.Enqueue(&m)
}
// 发送广播
func(this *Sender) SendNotice(msg interface{}){
	this.pool.Foreach(func(uid string){
		m:=queue.Message{
			MsgId:fmt.Sprintf("2%s",Guid()),
			MsgType:queue.Message_ToAll,
			Content:msg,
			To:uid,
		}
		this.queue.Enqueue(&m)
	})
}
// 客户端新上线处理
func(this *Sender) UpOnline(userID string){
	fmt.Printf("%s up line\n",userID)
	if this.Container==nil{
		return
	}	
	msgs,err:=this.Container.GetUserMsg(userID)
	if err!=nil{
		return
	}
	for _,v:=range msgs{
		this.queue.Enqueue(&v.Message)
	}
}


func(this *Sender) send(msg *queue.Message) error{
	msg.TrySendTimes++
	msg.SendTime=time.Now()

	c,err:=this.pool.Get(msg.To)
	if(err==nil){
		_,err=this.tcpProxy.Write(c,msg.Content)
	}
	return err	
}
// 消息发送失败处理
func (this *Sender)msgSendFailed(msg *queue.Message){
	//推送消息不处理
	if(msg.MsgType==queue.Message_ToAll){
		return
	}
	//在重试次数范围内，消息重新入队
	if this.ResendTimes >= msg.TrySendTimes{
		this.queue.Enqueue(msg)
		return
	}
	if this.Container==nil{
		return
	}
	m:=&persistence.OfflineMsg{
		Message:*msg,
	}
	this.Container.Put(m)
	
}