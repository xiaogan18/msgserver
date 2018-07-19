package msgserver
import(
	"github.com/xiaogan18/msgserver/pool"
	"github.com/xiaogan18/msgserver/serialize"
	"github.com/xiaogan18/msgserver/protocol"
	"github.com/xiaogan18/msgserver/queue"
	"fmt"
)
// 使用默认参数创建listener/sender
func NewDefaultServer(onSSL bool) (sder *SenderScheduler,lster *Listener,err error){
	return NewServer("default","default","default","default",onSSL)
}

func NewServer(poolType,queueType,serializer,protocolType string,OnSSL bool) (sder *SenderScheduler,lster *Listener,err error){
	defer func(){
		if e:=recover();e!=nil{
			sder=nil
			lster=nil
			err=fmt.Errorf("create msg server error:%s",e)
		}
	}()
	pl:=pool.CreatePool(poolType)
	qu:=queue.CreateQueue(queueType)
	// 初始化消息收发代理
	proxy:=&TcpProxy{
		Seri:serialize.CreateSerializer(serializer),
		Proto:protocol.CreatePro(protocolType),
		IsOnSSL:OnSSL,  //是否开启SSL
	}
	//初始化监听器
	lster=NewListener(pl,proxy)
	//初始化消息发送器
	sder=NewSender(pl,qu,proxy)
	//设置上线消息触发
	lster.onNewOnline=sder.UpOnline
	return
}