package main

import (
	"fmt"
	"sync"
)

func PrintEvent(x int,wg *sync.WaitGroup){
	if x %2 ==0{
		fmt.Printf("%d is even \n",x)
	}
	wg.Done()
}


func main() {

	var wg sync.WaitGroup
	for i:=0;i <= 1000;i++ {
		wg.Add(1)
		go PrintEvent(i,&wg)
	}
	wg.Wait()

}
