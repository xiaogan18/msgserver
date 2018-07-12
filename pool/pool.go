package pool
import(
	"net"
)

type Pool interface{
	Put(id string, conn net.Conn) error
	Get(id string) (net.Conn,error)
	Foreach(callback func(string))
	Clear()
	Count() int
}
func CreatePool(t string) (p Pool){
	switch(t){
	default:
		temp:=new(PoolMemory)
		temp.Init(60)
		p=temp
		break
	}
	return
}