package pool
import(
	"net"
)

type Pool interface{
	Put(id string, conn net.Conn) error
	Get(id string) (net.Conn,error)
	Foreach(callback func(net.Conn))
	Clear()
	Count() int
}