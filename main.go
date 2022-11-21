package main

import(
	"fmt"
	"sync"
	"db/src"
	"time"
)
var bit bool=false
func Put(primary string,wg *sync.WaitGroup, secondary []string,key string, value string){
	fmt.Println("Writing the Data to  server: key - ",key," value - ",value)
	length:=len(secondary)
	f:=2*length/3-1
	fmt.Println(f)

	document,_:=src.OpenAndLoad(primary)
	document.Set(key,value)
	document.Close()
	// out:=document.Get(key)
	fmt.Println("Updated server: ",primary,"data file: ",document)
	
	for i:=0;i<len(secondary);i++{
		time.Sleep(2 * time.Second)
		document,_=src.OpenAndLoad(secondary[i])
		document.Set(key,value)
		document.Close()
		// out:=document.Get(key)
		fmt.Println("Updated server: ",secondary[i],"data file: ",document)
		if(i==f){
			fmt.Println("Write operation completed...")
		}
		if(i>f){
			bit=true
		}
	}
	defer wg.Done()

}
func read(secondary []string,wg *sync.WaitGroup, key string){
	for{
		if bit{
			numbers := []int{0,1,2,3,4,5,6} 

			for i:=range numbers{
				document,_:=src.OpenAndLoad(secondary[7])
				k:=document.Get(key)
				// fmt.Printf("%T",k)
				fmt.Println(i,"server :",secondary[7]," output: ",key,"value: ",k)
				time.Sleep(1 * time.Second)

			}
			break
		} 
	}
	bit=false
	defer wg.Done()
	// val := k.kv[key]

}

func main()  {
	primary:="1001"
	secondary:=[]string{"1002","1003","1004","1005","1006","1007","1008","1009","1010"}
	src.CreateServer(primary)
	for i:=0;i<len(secondary);i++{
		src.CreateServer(secondary[i])
	}
	
	var wg sync.WaitGroup
	wg.Add(4)


	
	go Put(primary,&wg,secondary,"john","23")
	go read(secondary,&wg,"john")
	fmt.Println("----------------------------")
	go Put(primary,&wg,secondary,"jack","45")
	read(secondary,&wg,"jack")


	wg.Wait()
	fmt.Println("Done!")

	

}