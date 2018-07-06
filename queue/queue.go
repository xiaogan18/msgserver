package queue

type Queue interface{
	// 入队
	Enqueue(msg *Message) error
	// 出队
	Dequeue() (*Message,error)
}