package msgserver
import(
	"net"
)
type Filter interface{
	OnFilter(net.Conn) bool
}