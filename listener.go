package msgserver
import(
	"net"
	"fmt"
	"msgserver/pool"
)
type Lister struct{
	net.Listener
	pool pool.Pool
	Filter Filter
	//有新连接时发生
	OnAcceptCallback func()
	//身份验证
	OnAuthentication func(string) (string,error)
	//tcp代理类
	tcpProxy *TcpProxy
}
func(this *Lister) Init(pool pool.Pool,proxy *TcpProxy) error{
	this.pool=pool
	this.tcpProxy=proxy
	return nil
}
// 开启监听（阻塞）
func(this *Lister) Listen(address string) error{
	lster,err:=net.Listen("tcp",address)
	if(err!=nil){
		return err
	}

	for{
		conn,err:=lster.Accept()
		if(err==nil){
			go func(){
				if(this.OnAcceptCallback!=nil){
					this.OnAcceptCallback()
				}
				this.handler(conn)
			}()
		}
	}
	return nil
}

func (this *Lister)handler(conn net.Conn){
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
		connID=fmt.Sprintf("conn%d",this.pool.Count()+1)
	}
	this.pool.Put(connID,conn)
}
func (this *Lister) closeConn(conn net.Conn,err error){
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