package serialize

type Serialize interface{
	ToBytes(content interface{}) ([]byte,error)
	ToInterface(bytes []byte,obj interface{}) error
}