package queue

type Queue interface{
	// 入队
	Enqueue(msg *Message) error
	// 出队
	Dequeue() (*Message,error)
	// 消息总数
	Count() int
}

func CreateQueue(t string) (q Queue){
	switch(t){
	default:
		q=new(QueueMemory)
		break
	}
	return
}