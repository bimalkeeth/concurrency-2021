package main

import "fmt"

func main() {
	ch:=make(chan string)
	go func(){
		ch<-"Hello Channel"

	}()
	message:= <-ch
	fmt.Printf("Message %s",message)

}
