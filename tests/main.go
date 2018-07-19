package main
import(
	"msgserver"
	"bufio"
	"os"
	"fmt"
	"strings"
	"net"
	"errors"
	"github.com/xiaogan18/msgserver/persistence"
)
type testFilter struct{
}
func (this testFilter) OnFilter(conn net.Conn) bool{
	fmt.Println(conn.RemoteAddr().String())
	return true
}
func main(){
	sdr,lster,err:=msgserver.NewDefaultServer(false)
	if err!=nil{
		fmt.Println(err)
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
	// 开启协程 监听主机3365端口
	go func(){
		apier:=msgserver.Api{
			Sender:&sdr.Sender,
		}
		// 开启http api
		apier.Listen("127.0.0.1:3365","/msg/sender")
	}()
	// 设置离线消息容器
	sdr.Container=persistence.CreateMsgContainer("default")
	// 消息发送处理最大并行数设置
	sdr.MaxParallel=10
	//最小并行数设置
	sdr.MinParallel=1
	//队列中消息堆积阈值设置
	sdr.QueueBufferLen=1000
	// 开始sender
	sdr.BeginSender()
	for{
		// 从标准输入读取字符串，以\n为分割
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if(err==nil){
			text= strings.Replace(text,"\r\n","",1)
			if(text=="count"){
				fmt.Println(lster.OnlineCount())
			}else{
				sdr.SendNotice(text)
			}
		}
	}
}