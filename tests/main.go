package main
import(
	"msgserver"
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
	sdr,lster,err:=msgserver.NewDefaultServer()
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