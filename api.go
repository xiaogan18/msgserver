package msgserver
import(
	"net/http"
	"strings"
	"encoding/json"
	"fmt"
	"io"
)
type Api struct{
	Sender *Sender
}
func(this *Api) Listen(addr,path string){
	h:=&httpHandler{
		sender:this.Sender,
	}
	h.router=make(map [string]func(*Sender,http.ResponseWriter,*http.Request),0)
	h.router[path]=msgSender
	fmt.Printf("api listen address %s,path '%s'[post]\n",addr,path)
	err:=http.ListenAndServe(addr,h)
	if err!=nil{
		panic(err)
	}
}
// 消息发送 http api
func msgSender(sender *Sender,res http.ResponseWriter,req *http.Request){
	if strings.ToUpper(req.Method)!="POST"{
		return
	}
	var success=true
	defer func(){
		if(success){
			res.Write([]byte("1"))
		}else{
			res.Write([]byte("0"))
		}
	}()
	reader:= req.Body
	buffer:=make([]byte,256)
	le,err:=reader.Read(buffer)
	if err!=nil && err!=io.EOF{
		success=false
		return
	}
	buffer=buffer[:le]
	type msgContent struct{
		To string
		Content string
	}
	msg:=new(msgContent)
	if json.Unmarshal(buffer,msg)!=nil{
		success=false
		return
	}
	if msg.To==""{
		sender.SendNotice(msg.Content)
	}else{
		sender.SendMessage(msg.Content,msg.To)
	}
}

// http handler
type httpHandler struct{
	sender *Sender
	router map [string]func(*Sender,http.ResponseWriter,*http.Request)
}
func(this *httpHandler) ServeHTTP(res http.ResponseWriter, req *http.Request){
	path:=req.URL.Path
	for k:=range this.router{
		if k==path{
			this.router[k](this.sender, res,req)
			break
		}
	}
}