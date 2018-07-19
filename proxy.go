package msgserver
import(
	"net"
	"github.com/xiaogan18/msgserver/protocol"
	"github.com/xiaogan18/msgserver/serialize"
)
type TcpProxy struct{
	Proto protocol.Protocol
	Seri serialize.Serialize
	IsOnSSL bool
}
// 发送消息
func(this *TcpProxy) Write(conn net.Conn,data interface{}) (int,error){
	var err error
	var b []byte
	if b,err=this.Seri.ToBytes(data);err==nil{
		if this.IsOnSSL{  //是否开启SSL加密传输
			if b,err=Encrypt(b);err==nil{
				b=this.Proto.Packet(b)
			}
		}else{
			b=this.Proto.Packet(b)
		}
		return conn.Write(b)
	}
	return 0,err
}
// 接收消息
func(this *TcpProxy) Read(conn net.Conn) (chan []byte,error){
	var err error
	var b []byte=make([]byte,128)
	var i int
	if i,err=conn.Read(b);err==nil{
		b=b[:i]
		rdChan:=make(chan []byte)
		if this.IsOnSSL{   //是否开启SSL加密传输
			format:=func(data []byte) []byte{
				if r,err:=Decrypt(data);err!=nil{
					return nil
				}else{
					return r
				}
			}
			this.Proto.Unpack(b,rdChan,format)
		}else{
			this.Proto.Unpack(b,rdChan,nil)
		}
		return rdChan,nil
	}
	return nil,err
}
// 反序列化
func(this *TcpProxy) DeSerialize(data []byte,v interface{}) error{
	return this.Seri.ToInterface(data,v)
}
// 开始SSL握手
func(this *TcpProxy) SSL(conn net.Conn) error{
	if(this.IsOnSSL){
		return handshake(conn)
	}
	return nil
}