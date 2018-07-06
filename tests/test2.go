package main
import(
	"fmt"
)
var chans=make([]chan int ,4)
var files=make([]string,4)
// func write(i int){
// 	switch{
// 	case <-chans[0]:
// 		input(1)
// 	case <-chans[1]:
// 		input(2)
// 	case <-chans[2]:
// 		input(3)
// 	case <-chans[3]:
// 		input(4)
// 	}
// }
func input(f_id,i int){
	files[f_id]+=fmt.Sprintf("%d",i)
	//fmt.Println(files[f_id])
}
func main(){
	for index:=0;index<4;index++ {
		chans[index]=make(chan int)
		go func(i int){
			for{
				//fmt.Printf("go %d print\n",i)
				f:=<-chans[i]
				input(f,i+1)
			}
		}(index)
	}
	fmt.Println("go run..")
	for times:=0;times<4;times++ {
		for i:=0;i<4;i++ {
			i=times+i
			if i>=4{
				i=i-4
			}
			//fmt.Printf("file %d write\n",times)
			chans[i]<-times
		}
	}
	fmt.Println(files)
}