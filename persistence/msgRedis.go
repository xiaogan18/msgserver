package persistence
import(
	"github.com/gomodule/redigo/redis"
	"github.com/xiaogan18/msgserver/serialize"
	"log"
	"fmt"
)
const(
	userMsgPrefix="u_m_"
)
type RedisOptions struct{
	Network string
	Address string
	Password string
}
type RedisContainer struct{
	option *RedisOptions
	serializer serialize.Serialize
}
func(this *RedisContainer) try(){
	c,err:=this.dial()
	if err!=nil{
		panic(fmt.Sprintf("redis conn error:%s",err))
	}
	defer func(){
		c.Close()
	}()
}
func(this *RedisContainer) dial() (redis.Conn,error){
	var op=make([]redis.DialOption,0)
	if this.option.Password!=""{
		op=append(op,redis.DialPassword(this.option.Password))
	}
	c,err:=redis.Dial(this.option.Network,this.option.Address,op...)
	return c,err
}
func(this *RedisContainer) Get(id string) (*OfflineMsg,error){
	c,err:=this.dial()
	if err!=nil{
		return nil,err
	}
	defer func(){
		c.Close()
	}()
	// 从Redis中取出
	jsonStr,err:=redis.String(c.Do("GET",id))
	if err!=nil{
		return nil,err
	}
	// 反序列化
	msg:=new(OfflineMsg)
	if err=this.serializer.ToInterface([]byte(jsonStr),msg);err!=nil{
		return nil,err
	}
	//删除redis中的消息
	c.Do("DEL",id)
	return msg,nil
}
func(this *RedisContainer) GetUserMsg(userID string) ([]*OfflineMsg,error){
	c,err:=this.dial()
	if err!=nil{
		return nil,err
	}
	defer func(){
		c.Close()
	}()
	// 取出user下msg id
	msgs,err:= redis.Values(c.Do("SMEMBERS",userMsgPrefix + userID))
	if err!=nil{
		return nil,err
	}
	result:=make([]*OfflineMsg,0)
	// 循环取msg
	for _,k:=range msgs{
		id,_:=redis.String(k,nil)
		jsonStr,err:=redis.String(c.Do("GET",id))
		if err==nil{
			msg:=new(OfflineMsg)
			if err=this.serializer.ToInterface([]byte(jsonStr),msg);err==nil{
				result=append(result,msg)
			}
			// 删除消息
			c.Do("DEL",id)
		}
	}
	// 删除用户离线消息列表
	c.Do("DEL",userMsgPrefix + userID)
	return result,nil
}
func(this *RedisContainer) Put(msg *OfflineMsg){
	c,err:=this.dial()
	if err!=nil{
		log.Printf("offline msg input failed:%s",err)
	}
	defer func(){
		c.Close()
	}()
	json,_:=this.serializer.ToBytes(msg)
	c.Do("SET",msg.MsgId,string(json))
	c.Do("SADD",userMsgPrefix+msg.To,msg.MsgId)
}