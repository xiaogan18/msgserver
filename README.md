# msgserver
## 说明
一种基于socket的消息推送服务端，服务端保持与客户端tcp长连接，实现消息的推送或提醒。<br>
目前支持消息模式有两种：1.全员广播；2.特定用户推送<br>
现在还只能支持单台服务器部署，后期将提供分布式部署方案<br>
<br>
## 快速使用
### 安装
	go get github.com/xiaogan18/msgserver
### 创建hello.go
```go
package main
import(
	"github.com/xiaogan18/msgserver"
	"fmt"
	"bufio"
	"strings"
	"os"
)
func main(){
	sdr,lster,err:=msgserver.NewDefaultServer(false)  //不开启ssl加密传输
	if err!=nil{
		fmt.Println(err)
	}
	//开启一个协程 监听主机3366端口
	go func(){
		lster.Listen("127.0.0.1:3366")
	}()
	sdr.BeginSender()

	for{
		// 从标准输入读取字符串，以\n为分割
		fmt.Println("input a msg:")
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
```
### 运行
	go run hello.go
	使用客户端tcp连接到127.0.0.1:3366就可以正常收到消息了
## 详细描述
### 主要流程
![image](https://github.com/xiaogan18/msgserver/blob/master/.github/主流程图.png)
### 带参数构建服务
```go
//函数签名
func NewServer(poolType,queueType,serializer,protocolType string,OnSSL bool) (*SenderScheduler,*Listener,error)
```
	poolType:连接池类型（默认'default'）
	queueType:队列类型（默认'default'）
	serializer:序列化器类型（默认'json'）
	protocolType:协议类型（默认'default'）
	OnSSL:是否开启SSL加密
```go
//带参数构建服务实体
sdr,lster,err:=msgserver.NewServer("default","default","json","default",true)
```
### SSL握手
	如果你开启了SSL，建立连接时将触发三次握手密钥交换，且之后的所有消息将使用加密传输。 <br>
	密钥交换使用RSA非对称加密算法，密钥长度1024，填充方式PKCS1，密文过长使用分段加密方式
	消息加密使用AES对称加密算法，密钥长度16，模式ECB，填充方式PKCS7
	握手流程如下：
![image](https://github.com/xiaogan18/msgserver/blob/master/.github/SSL流程.png)
### 过滤器
	如果你想过滤某些连接，应用黑/白名单，可以使用过滤器
```go
// 定义struct满足msgserver.Filter接口
type testFilter struct{
}
func (this testFilter) OnFilter(conn net.Conn) bool{
	fmt.Println(conn.RemoteAddr().String())
	return true
}
```
```go
// 将过滤器应用到Lisenter
lster.Filter=new(testFilter)
```
### 身份验证
	如果你需要验证客户端的身份，可以定义身份验证函数。
	其中string参数是建立连接后客户端发送来的一段字符串,验证失败服务端会主动断开连接
```go
lster.OnAuthentication=func(data string) (string,error){
	user:=strings.Split(data,";")
	if len(user)==2 && user[1]=="123456"{
		return user[0],nil
	}else{
		return "",errors.New("password is wrong")
	}
}
```
### 离线消息持久化
	如果你需要保存用户离线后的消息，可以使用持久化器。
	用户下次上线后会立即取出消息发送（只会对调用Sender.SendMessage()的消息进行持久化）
```go
// 保存在内存
sdr.Container=persistence.CreateMsgContainer("default")  
// 保存在redis
sdr.Container=persistence.CreateMsgContainer("redis",
	&persistence.RedisOptions{Network:"tcp",Address:"127.0.0.1:6379",})
```
### HTTP API接口调用
	如果你想通过http来控制消息发送，可以使用api监听
```go
go func(){
	apier:=msgserver.Api{
		Sender:&sdr.Sender,
	}
	// 开启http api
	apier.Listen("127.0.0.1:3365","/msg/sender")  //监听主机端口，服务目录“/msg/sender” Method="POST"
}()
```
http post body(To为空表示发送全员通知):
{
  	"To":"xiaogan18",
  	"Content":"hello world"
}
### Sender调度器参数
```go
// 消息发送重试次数
sdr.ResendTimes=2
// 重试间隔（毫秒）
sdr.ResendInterval=10*1000
// 消息发送处理最大并行数设置
sdr.MaxParallel=10
//最小并行数设置
sdr.MinParallel=1
//队列中消息堆积阈值设置
sdr.QueueBufferLen=1000
```
### 开始服务
	如果你完成了所有配置，就可以按下面的方式开始服务了
```go
//listener开始监听
go func(){
	lster.Listen("127.0.0.1:3366")
}()
// 开始sender
sdr.BeginSender()
```
