package main

import (
	"fmt"
	"math/rand"
	"time"
)

func cakeMaker(kind string,number int,to chan<- string){
    rand.Seed(time.Now().Unix())
    for i:=0;i<number;i++{
    	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
    	to <-kind
	}
	close(to)
}

func main() {

	chocolateChan:=make(chan string,8)
	redVelvetChan:=make(chan string,8)

	go cakeMaker("chocolate",10,chocolateChan)
	go cakeMaker("red velvet",8,redVelvetChan)

	moreChocolate :=true
	moreVelvet :=true
	var cake string

	for moreChocolate || moreVelvet {
		select {
		    case cake,moreChocolate=<-chocolateChan:{
                 if moreChocolate{
                 	fmt.Printf("Got a cake from the first factory: %s\n",cake)
				 }
			}
			case cake,moreVelvet=<-redVelvetChan:{
				if moreVelvet{
					fmt.Printf("Got a cake from the second factory: %s\n",cake)
				}
			}
			case <-time.After(3 * time.Second):
				fmt.Println("Time out")
				return
		}
	}




}
