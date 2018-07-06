package pool
import(
	"net"
	//"fmt"
	"time"
)
// 开始自动清理无效连接
func beginConnGC(connList *map[string]net.Conn,interval int64){
	go func(){
		for{
			list:=*connList
			downConns:=make([]string,0)
			for key:=range list{
				conn:=list[key]
				if _,err:=conn.Write([]byte{0});err!=nil{
					downConns= append(downConns,key)
				}
			}
			//移除无效连接conn
			//fmt.Println("begin conn gc")
			_lock.Lock()
			for _,k:=range downConns{
			//	fmt.Println("remove a conn "+k)
				delete(list,k)
			}
			_lock.Unlock()
			time.Sleep(time.Duration(interval) * time.Millisecond)  //休眠
		}
	}()
}