package serialize
import(
	"encoding/json"
)

type JsonSerialize struct{

}

func(this *JsonSerialize) ToBytes(content interface{}) ([]byte,error){
	//s:= content.(string)
	data,err:=json.Marshal(content)
	return data,err
}
func(this *JsonSerialize) ToInterface(bytes []byte,obj interface{}) error{

	return json.Unmarshal(bytes,obj)
}