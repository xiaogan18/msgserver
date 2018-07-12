package queue
import(
	"errors"
	"container/list"
	"sync"
)
var(
	_lock=new (sync.Mutex)    //互斥锁
	Error_QueueIsEmpty=errors.New("queue is already empty")   //队列为空
)
  
type QueueMemory struct{
	queue *list.List
	//类型 0:先入先出
	queueType int

}
//入队
func(this *QueueMemory) Enqueue(msg *Message) error{
	if this.queue==nil{
		this.queue=list.New()
	}
	this.queue.PushBack(msg)
	return nil
}
//出队
func(this *QueueMemory) Dequeue() (*Message,error){
	if this.queue==nil || this.queue.Len()==0{
		return nil,Error_QueueIsEmpty
	}
	var msg *Message
	//判断队列类型
	switch(this.queueType){
	case 0:  //先入先出
		ele:=this.queue.Front()
		_lock.Lock()   //修改时使用 互斥锁
		msg= this.queue.Remove(ele).(*Message)
		_lock.Unlock()
		break
	}
	return msg,nil
}
func(this *QueueMemory) Count() int{
	if this.queue==nil{
		return 0
	}
	return this.queue.Len()
}