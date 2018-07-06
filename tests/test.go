package main
import(
	"fmt"
	"sync"
	"runtime"
)
var ch1 =make(chan int)
//var ch2=make(chan int)
const(
	nums="123456789"
	letters="abcdefghi"
)
func main(){
	runtime.GOMAXPROCS(1)
	wd:=sync.WaitGroup{}
	wd.Add(len(nums)+len(letters))
	fun1:=func(){
		for _,v:=range nums{
			// if(i%2==0 && i>1){
			// 	ch2<-1
			// 	<-ch1
			// }
			// if(i==0){
			// 	ch1<-0
			// }
			// else if(i%2==0){
			// 	ch1<-1
			// }
			fmt.Print(string(v))
			ch1<-1
			wd.Done()
		}
		ch1<-1
	}
	fun2:=func(){
		// if c:=<-ch1;c==0{
		// 	<-ch1
		// }
		<-ch1
		for _,v:=range letters{
			// if(i%2==0 && i>1){
			// 	ch1<-1
			// 	<-ch2
			// }
			// if(i%2==0 && i>1){
			// 	//<-ch1
			// }
			fmt.Print(string(v))
			<-ch1
			wd.Done()
		}
	}
	go fun1()
	go fun2()
	wd.Wait()
}