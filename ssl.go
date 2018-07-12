package msgserver
import(
	"crypto/rsa"
	"crypto/x509"
	"crypto/rand"
	"crypto/aes"
	"crypto/sha256"
	"time"
	"net"
	"fmt"
)
var(
	privateKey *rsa.PrivateKey
	rsa_privateKey []byte
	rsa_publicKey []byte
	secretKey []byte
)
func init(){
	var err error
	privateKey,err=rsa.GenerateKey(rand.Reader,1024)
	if err!=nil{
		panic("failed to create rsa secret key!")
	}
	rsa_privateKey= x509.MarshalPKCS1PrivateKey(privateKey)
	rsa_publicKey= x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	//用当前时间sha256生成一个秘钥
	now:=time.Now().Unix()
	hash:=sha256.Sum256([]byte(fmt.Sprintf("%d",now)))
	secretKey=hash[:16]
}
// 开始三次握手
func handshake(conn net.Conn) error{
	// 发送公钥
	if _,err:=conn.Write(rsa_publicKey);err!=nil{
		return err
	}
	// 接收客户端密文
	text:=make([]byte,128)
	i,err:=conn.Read(text)
	if(err!=nil){
		return err
	}
	text=text[:i]
	// 解密得到客户端publicKey
	clientKey,err:= rsa.DecryptPKCS1v15(rand.Reader,privateKey,text)
	if(err!=nil){
		return err
	}
	cKey,err:= x509.ParsePKCS1PublicKey(clientKey)
	if err!=nil{
		return err
	}
	// 使用客户端publicKey加密 对称密钥secretKey
	lastText,err:=rsa.EncryptPKCS1v15(rand.Reader,cKey,secretKey)
	if err!=nil{
		return err
	}
	// 将secretKey加密后的密文发送到客户端
	if _,err=conn.Write(lastText);err!=nil{
		return err
	}
	return nil
}
// 加密
func Encrypt(data []byte) ([]byte,error){
	b,err:= aes.NewCipher(secretKey)
	if(err!=nil){
		return nil,err
	}
	secretText:=make([]byte,0)
	b.Encrypt(secretText,data)
	return secretText,nil
}
 // 解密
func Decrypt(data []byte) ([]byte,error){
	b,err:= aes.NewCipher(secretKey)
	if(err!=nil){
		return nil,err
	}
	text:=make([]byte,0)
	b.Decrypt(text,data)
	return text,nil
}