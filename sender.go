package msgserver
import(
	"msgserver/queue"
	"msgserver/pool"
	"net"
	"time"
)

type Sender struct{
	pool pool.Pool
	queue queue.Queue
	//tcp代理类
	tcpProxy *TcpProxy
	FailedCallback func(error)   //消息发送失败时回调
}
//  初始化
func(this *Sender) Init(p pool.Pool,q queue.Queue,proxy *TcpProxy){
	this.pool=p
	this.queue=q
	this.tcpProxy=proxy
	//开启线程 取队列消息
	go func(){
		for{
			m,err:=this.queue.Dequeue()
			if(err==nil){
				if err:=this.send(m);err!=nil{
					if(this.FailedCallback!=nil){  //发送失败回调
						this.FailedCallback(err)
					}
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
func(this *Sender) send(msg *queue.Message) error{
	var err error
	switch(msg.MsgType){
		case queue.Message_ToOne:  //发送单条消息
			var c net.Conn
			c,err=this.pool.Get(msg.To)
			if(err!=nil){
				_,err=this.tcpProxy.Write(c,msg.Content)
			}
			break
		case queue.Message_ToAll:  //群发消息
			this.pool.Foreach(func(c net.Conn){
				this.tcpProxy.Write(c,msg.Content)
			})
			break
	}
	return err	
}
