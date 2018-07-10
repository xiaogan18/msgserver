package msgserver
import(
	"msgserver/queue"
	"msgserver/pool"
	"msgserver/persistence"
	"net"
	"time"
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
//  初始化
func(this *Sender) Init(p pool.Pool,q queue.Queue,proxy *TcpProxy){
	this.pool=p
	this.queue=q
	this.tcpProxy=proxy
	//开启协程 取队列消息
	go func(){
		for{
			m,err:=this.queue.Dequeue()
			if(err==nil){
				//判断是否重发消息
				if m.TrySendTimes>0{
					sTime:=m.SendTime.Add(time.Duration(this.ResendInterval)*time.Millisecond)
					if(time.Now().Before(sTime)){
						this.queue.Enqueue(m)   //未到重发时间，重新入队
						continue
					}
				}
				if err:=this.send(m);err!=nil{
					if(this.FailedCallback!=nil){  //发送失败回调
						this.FailedCallback(err)
					}
					//失败处理
					this.msgSendFailed(m)
				}
			}else{
				time.Sleep(1*time.Second)   //队列已空，休息
			}
		}
	}()
	
}
// 发送消息
func(this *Sender) SendMessage(msg interface{},user string) error{
	m:=queue.Message{
		MsgId:"1",
		MsgType:queue.Message_ToOne,
		Content:msg,
		To:user,
	}
	return this.queue.Enqueue(&m)
}
// 发送广播
func(this *Sender) SendNotice(msg interface{}) error{
	m:=queue.Message{
		MsgId:"2",
		MsgType:queue.Message_ToAll,
		Content:msg,
	}
	return this.queue.Enqueue(&m)
}
// 客户端新上线处理
func(this *Sender) UpOnline(userID string){
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
	var err error
	switch(msg.MsgType){
		case queue.Message_ToOne:  //发送单条消息
			var c net.Conn
			c,err=this.pool.Get(msg.To)
			if(err!=nil){
				_,err=this.tcpProxy.Write(c,msg.Content)
			}
			break
		case queue.Message_ToAll:  //推送消息
			this.pool.Foreach(func(c net.Conn){
				this.tcpProxy.Write(c,msg.Content)
			})
			break
	}
	return err	
}
// 消息发送失败处理
func (this *Sender)msgSendFailed(msg *queue.Message){
	//推送消息不处理
	if(msg.MsgType==queue.Message_ToAll){
		return
	}
	//再重试次数范围内，消息重新入队
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