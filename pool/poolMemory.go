package pool
import(
	"net"
	"errors"
	"sync"
)
var(
	_lock=new (sync.Mutex)    //互斥锁
	Err_NotFoundConn=errors.New("conn is not found")
)

type PoolMemory struct{
	pool map[string]net.Conn
}
// 初始化连接池
func(this *PoolMemory) Init(gcInterval int64){
	this.pool=make(map[string]net.Conn)
	beginConnGC(&this.pool,gcInterval*1000)
}

func(this *PoolMemory) Put(id string, conn net.Conn) error{
	_lock.Lock()
	this.pool[id]=conn
	_lock.Unlock()
	return nil
}
func(this *PoolMemory) Get(id string) (net.Conn,error){
	v,ok:=this.pool[id]
	if(!ok){
		return nil,Err_NotFoundConn
	}
	return v,nil
}
func(this *PoolMemory) Foreach(callback func(net.Conn)){
	if(this.pool!=nil && len(this.pool)>0){
		for _,v:= range this.pool{
			callback(v)
		}
	}
}
func(this *PoolMemory) Clear(){
	defer func(){
		this.pool=nil
	}()
	this.Foreach(func(conn net.Conn){
		conn.Close()
	})
}
func(this *PoolMemory) Count() int{
	if(this.pool==nil){
		return 0
	}
	return len(this.pool)
}