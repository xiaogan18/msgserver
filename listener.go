package msgserver
import(
	"net"
	"fmt"
	"msgserver/pool"
)
type Listener struct{
	net.Listener
	pool pool.Pool
	//tcp代理类
	tcpProxy *TcpProxy
	//有客户端上线时发生
	onNewOnline func(string)
	// 过滤器
	Filter Filter
	//身份验证
	OnAuthentication func(string) (string,error)
}
// 初始化一个监听器
func NewListener(pool pool.Pool,proxy *TcpProxy) *Listener{
	this:=&Listener{}
	this.pool=pool
	this.tcpProxy=proxy
	return this
}
// 开启监听（阻塞）
func(this *Listener) Listen(address string) error{
	lster,err:=net.Listen("tcp",address)
	if(err!=nil){
		return err
	}
	fmt.Printf("listen address %s ...\n",address)
	for{
		conn,err:=lster.Accept()
		if(err==nil){
			go func(){
				this.handler(conn)
			}()
		}
	}
}
func(this *Listener) OnlineCount() int{
	return this.pool.Count()
}

func (this *Listener)handler(conn net.Conn){
	//过滤器过滤连接
	if(this.Filter!=nil){
		if b:=this.Filter.OnFilter(conn);!b{
			this.closeConn(conn,nil)
			return
		}
	}
	//SSL
	if err:=this.tcpProxy.SSL(conn);err!=nil{
		this.closeConn(conn,err)
		return 
	}
	var connID string
	//身份验证
	if(this.OnAuthentication!=nil){
		buffer,err:=this.tcpProxy.Read(conn)
		if err!=nil{
			this.closeConn(conn,err)
			return
		}
		readStr:=string(<-buffer)
		if connID,err=this.OnAuthentication(readStr);err!=nil{
			this.closeConn(conn,err)
			return
		}
	}
	if connID==""{
		connID=fmt.Sprintf("conn%s",conn.RemoteAddr().String())
	}
	// 客户端上线触发
	if this.onNewOnline!=nil{
		this.onNewOnline(connID)
	}
	this.pool.Put(connID,conn)
}
func (this *Listener) closeConn(conn net.Conn,err error){
	defer func(){
		conn.Close()
	}()
	var errStr string
	if(err==nil){
		errStr="authorization validation failed"
	}else{
		errStr=fmt.Sprint(err)
	}
	conn.Write([]byte(errStr))
}