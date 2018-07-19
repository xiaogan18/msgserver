package msgserver
import(
	"net"
	"encoding/base64"
	"fmt"
)
var(
	rsa_privateKey []byte
	rsa_publicKey []byte
	secretKey []byte
)
func init(){
	// 创建rsa密钥对
	r:=rsa_encrypt{}
	if pr,pu,err:= r.newKey(1024);err==nil{
		rsa_privateKey=pr
		rsa_publicKey=pu
	}else{
		panic("failed to create rsa secret key!")
	}
	// 创建aes密钥
	secretKey=aes_encrypt{}.newKey(16)
}
// 开始三次握手  交换数据皆以base64格式
func handshake(conn net.Conn) error{
	// 发送公钥
	base64PubKey:=tobase64String(rsa_publicKey)
	if _,err:=conn.Write([]byte(base64PubKey));err!=nil{
		return err
	}
	// 接收客户端密文
	text:=make([]byte,512)
	i,err:=conn.Read(text)
	if(err!=nil){
		return err
	}
	text=text[:i]
	text,err=decodeBase64(text)
	if(err!=nil){
		return err
	}
	// 解密得到客户端publicKey
	clientKeyBase64,err:= rsa_encrypt{}.rsaDecrypt(text,rsa_privateKey)
	if err!=nil{
		return err
	}
	// publick key base64解码
	clientKey,err:=decodeBase64(clientKeyBase64)
	if err!=nil{
		return err
	}
	
	// 使用客户端publicKey加密 对称密钥secretKey
	lastText,err:=rsa_encrypt{}.rsaEncrypt(secretKey,clientKey)
	if err!=nil{
		return err
	}
	// 将secretKey加密后的密文base64发送到客户端
	if _,err=conn.Write([]byte(tobase64String(lastText)));err!=nil{
		return err
	}
	return nil
}
// 加密
func Encrypt(data []byte) ([]byte,error){
	return aes_encrypt{}.aesEncrypt(data,secretKey)
}
 // 解密
func Decrypt(data []byte) (result []byte,e error){
	defer func(){
		if err:=recover();err!=nil{
			e=fmt.Errorf("decrypt error:%s",err)
			fmt.Println(e)
		}
	}()
	result,e= aes_encrypt{}.aesDecrypt(data,secretKey)
	return
}

func tobase64String(data []byte) string{
	return base64.StdEncoding.EncodeToString(data)
}
func decodeBase64(data []byte) ([]byte,error){
	buffer:=make([]byte,512)
	i,err:=base64.StdEncoding.Decode(buffer,data)
	if err!=nil{
		return nil,err
	}
	return buffer[:i],nil
}