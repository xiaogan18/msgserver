package msgserver
import(
	"github.com/xiaogan18/msgserver/queue"
	"github.com/xiaogan18/msgserver/pool"
	"sync"
	"time"
	"errors"
	"fmt"
)
var(
	Err_TaskIsNotDefine=errors.New("task is not define")
)

//  初始化一个sender调度器
func NewSender(p pool.Pool,q queue.Queue,proxy *TcpProxy) *SenderScheduler{
	this:=&SenderScheduler{}
	this.pool=p
	this.queue=q
	this.tcpProxy=proxy
	//定义任务
	this.task= func(m *queue.Message){
		//判断是否重发消息
		if m.TrySendTimes>0{
			sTime:=m.SendTime.Add(time.Duration(this.ResendInterval)*time.Millisecond)
			if(time.Now().Before(sTime)){
				this.queue.Enqueue(m)   //未到重发时间，重新入队
				return
			}
		}
		if err:=this.send(m);err!=nil{
			if(this.FailedCallback!=nil){  //发送失败回调
				this.FailedCallback(err)
			}
			//失败处理
			this.msgSendFailed(m)
		}
	}
	return this
}
// 消息发送调度器
type SenderScheduler struct{
	Sender
	// 最大并行数
	MaxParallel int
	// 最小并行数
	MinParallel int
	// 队列堆积消息阈值，达到就增加并行
	QueueBufferLen int
	// 当前并行数
	parallelNum int
	task func(*queue.Message)
	sync.Mutex
}

func(this *SenderScheduler) BeginSender() error{
	if this.task==nil{
		return Err_TaskIsNotDefine
	}
	this.parallelNum=0
	if this.MaxParallel<=0{
		this.MaxParallel=10
	}
	if this.MinParallel<=0{
		this.MinParallel=1
	}
	if this.QueueBufferLen<=0{
		this.QueueBufferLen=1000
	}
	// 开启最小并行数量sender
	for i:=0;i<this.MinParallel;i++{
		this.newGoroutine()
	}
	//开启协程，监控队列中的消息有没有达到阈值
	go func(){
		for{
			if this.queue.Count() >= this.QueueBufferLen{
				this.newGoroutine()
			}
			time.Sleep(time.Second*10)
		}
	}()
	return nil
}
func(this *SenderScheduler) newGoroutine(){
	defer func(){
		this.Unlock()
	}()
	this.Lock()
	// 判断是否达到最大并行
	if this.parallelNum>=this.MaxParallel{
		return
	}else{
		this.parallelNum++
	}
	go func(this *SenderScheduler){
		defer func(){
			if err:=recover();err!=nil{
				fmt.Println(err)
			}
		}()
		for{
			m,err:=this.queue.Dequeue()  //取队列消息
			if err==queue.Error_QueueIsEmpty{
				this.Lock()
				b:=this.parallelNum<=this.MinParallel
				this.Unlock()
				//判断协程结束或阻塞
				if b{
					//已达到最小并行数量，休眠
					time.Sleep(time.Second*1)
				}else{
					break
				}
			}else if err==nil{
				//fmt.Println(*m)
				this.task(m)
			}
		}
	}(this)
	fmt.Println("a new sender groutine created")
}