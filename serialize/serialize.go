package serialize

type Serialize interface{
	ToBytes(content interface{}) ([]byte,error)
	ToInterface(bytes []byte,obj interface{}) error
}

func CreateSerializer(t string) (s Serialize){
	switch(t){
	default:
		s=new(JsonSerialize)
		break
	}
	return
}