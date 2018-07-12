package msgserver
import(
	"time"
	"math/rand"
	"crypto/sha256"
	"fmt"
)

func Guid() string{
	t:=time.Now().Unix()
	rm:=rand.Intn(1000000)
	str:=fmt.Sprintf("%d%d",t,rm)
	guByte:= sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x",guByte)
}