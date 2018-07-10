package protocol
import(
	"bytes"
)

// 传输协议
type Protocol interface{
	// 封包
	Packet(msg []byte) []byte
	// 解包
	Unpack(buffer []byte,readChan chan []byte)
}
func CreatePro(t string)(p Protocol){
	switch(t){
	default:
		p=new(CustomPro)
		break
	}
	return
}

type CustomPro struct{

}
const(
	proSplitChar='|'
)
func(this *CustomPro) Packet(msg []byte) []byte{
	return append(msg,byte(proSplitChar))
}
func(this *CustomPro) Unpack(buffer []byte,readChan chan []byte){
	go func(){
		startIndex:=0
		chars:=bytes.Runes(buffer)
		for i,c:=range chars{
			if c==proSplitChar{
				readChan<-buffer[startIndex:i]
				startIndex=i+1
			}
		}
		if startIndex >= len(buffer){
			readChan<-[]byte{}
		}else{
			readChan<-buffer[startIndex:]
		}
	}()
}