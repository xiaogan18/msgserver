package msgserver
import(
	"crypto/aes"
	"crypto/rsa"
	"crypto/x509"
	"time"
	"crypto/sha256"
	"bytes"
	"crypto/rand"
	"fmt"
)
type aes_encrypt struct{
}
// 随机生成一个key bits:16|24|32
func(this aes_encrypt) newKey(bits int) []byte{
	//用当前时间sha256生成一个AES秘钥
	now:=time.Now().Unix()
	hash:=sha256.Sum256([]byte(fmt.Sprintf("%d",now)))
	return hash[:bits]
}
// aes加密 模式：ECB 填充：PKCS7
func(this aes_encrypt) aesEncrypt(plantText,key []byte) ([]byte,error){
	block, err := aes.NewCipher(key) //选择加密算法 ECB
    if err != nil {
       return nil, err
    }
    plantText = this.pKCS7Padding(plantText, block.BlockSize())

    ciphertext := make([]byte, len(plantText))
	block.Encrypt(ciphertext,plantText)
	ciphertext=[]byte(tobase64String(ciphertext))  //密文以base64格式传输
    return ciphertext, nil
}
// aes解密 模式：ECB 填充：PKCS7
func(this aes_encrypt) aesDecrypt(ciphertext,key []byte) ([]byte, error) {
	var err error
	ciphertext,err=decodeBase64(ciphertext)  //密文以base64格式传输
	if err != nil {
		return nil, err
	 }
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes) //选择加密算法 ECB
	if err != nil {
	   return nil, err
	}
	plantText := make([]byte, len(ciphertext))
	block.Decrypt(plantText, ciphertext)
	plantText = this.pKCS7UnPadding(plantText, block.BlockSize())
	return plantText, nil
 }
 // pkcs7填充
func(this aes_encrypt) pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
 }
 // 去除pkcs7填充
 func(this aes_encrypt) pKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
 }

 type rsa_encrypt struct{

 }
 // 创建一对公私密钥 bits:1024|2048
 func(this rsa_encrypt) newKey(bits int) (privateKey,publicKey []byte,e error){
	p,err:=rsa.GenerateKey(rand.Reader,bits)
	if err!=nil{
		e=err
		return
	}
	privateKey= x509.MarshalPKCS1PrivateKey(p)
	publicKey,e= x509.MarshalPKIXPublicKey(&p.PublicKey)
	return
 }
 // rsa加密 PKCS1填充
 func(this rsa_encrypt) rsaEncrypt(plantText,publicKey []byte) ([]byte,error){
	cKey,err:= x509.ParsePKIXPublicKey(publicKey)
	if err!=nil{
		return nil,err
	}
	return rsa.EncryptPKCS1v15(rand.Reader,cKey.(*rsa.PublicKey),plantText)
 }
 // rsa解密 PKCS1填充
 func(this rsa_encrypt) rsaDecrypt(plantText,privateKey []byte) ([]byte,error){
	pKey,err:=x509.ParsePKCS1PrivateKey(privateKey)
	if(err!=nil){
		return nil,err
	}
	result:=make([]byte,0)
	length:=pKey.PublicKey.N.BitLen()/8  //单次密文长度
	fmt.Println(length)
	times:=len(plantText)/length
	for i:=0;i<times;i++{
		temp,err:=rsa.DecryptPKCS1v15(rand.Reader,pKey,plantText[i*length:(i+1)*length])
		if(err!=nil){
			return nil,err
		}
		result=append(result,temp...)
	}
	return result,nil
}