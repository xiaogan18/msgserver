package main
import(
	. "msgserver"
	"msgserver/pool"
	"msgserver/queue"
	"msgserver/serialize"
	"msgserver/protocol"
	"bufio"
	"os"
	"fmt"
	"strings"
	"net"
	"errors"
)
type testFilter struct{
}
func (this testFilter) OnFilter(conn net.Conn) bool{
	fmt.Println("filter return true")
	return true
}
func main(){
	// 初始化消息收发代理
	proxy:=&TcpProxy{
		Seri:&serialize.JsonSerialize{},
		Proto:&protocol.CustomPro{},
	}

	pl:=new(pool.PoolMemory)
	//初始化连接池 设置清理间隔为30s
	pl.Init(30)
	//初始化监听器
	lster:=new(Lister)
	lster.Init(pl,proxy)
	lster.OnAcceptCallback=func(){
		fmt.Println("a connetion accpet...")
	}
	//设置连接过滤器
	lster.Filter=new(testFilter)
	//设置身份验证方法
	lster.OnAuthentication=func(data string) (string,error){
		user:=strings.Split(data,";")
		if len(user)==2 && user[1]=="123456"{
			return user[0],nil
		}else{
			return "",errors.New("password is wrong")
		}
	}
	//开启一个协程 监听主机3366端口
	go func(){
		lster.Listen("127.0.0.1:3366")
	}()
	
	//初始化消息队列
	qu:=new(queue.QueueMemory)
	
	//初始化消息发送器
	sdr:=new(Sender)
	sdr.Init(pl,qu,proxy)
	for{
		// 从标准输入读取字符串，以\n为分割
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if(err==nil){
			text= strings.Replace(text,"\r\n","",1)
			if(text=="count"){
				fmt.Println(pl.Count())
			}else{
				sdr.SendNotice(text)
			}
		}
	}
}