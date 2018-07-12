package pool
import(
	"net"
	"errors"
	"sync"
)
var(
	Err_NotFoundConn=errors.New("conn is not found")
)

type PoolMemory struct{
	_lock sync.Mutex    //互斥锁
	pool map[string]net.Conn
}
// 初始化连接池
func(this *PoolMemory) Init(gcInterval int64){
	this.pool=make(map[string]net.Conn)
	this.beginConnGC(this.pool,gcInterval*1000)
}

func(this *PoolMemory) Put(id string, conn net.Conn) error{
	defer func(){
		this._lock.Unlock()
	}()
	this._lock.Lock()
	this.pool[id]=conn
	return nil
}
func(this *PoolMemory) Get(id string) (net.Conn,error){
	defer func(){
		this._lock.Unlock()
	}()
	this._lock.Lock()
	v,ok:=this.pool[id]
	if(!ok){
		return nil,Err_NotFoundConn
	}
	return v,nil
}
func(this *PoolMemory) Foreach(callback func(string)){
	if(this.pool!=nil && len(this.pool)>0){
		for v:= range this.pool{
			callback(v)
		}
	}
}
func(this *PoolMemory) Clear(){
	defer func(){
		this._lock.Unlock()
		this.pool=nil
	}()
	this._lock.Lock()
	for k:=range this.pool{
		this.pool[k].Close()
	}
}
func(this *PoolMemory) Count() int{
	if(this.pool==nil){
		return 0
	}
	defer func(){
		this._lock.Unlock()
	}()
	this._lock.Lock()
	return len(this.pool)
}