package main

import (
	"fmt"
	"net/http"
	"sync"
)

func WebGetWorker(in <-chan string,wg *sync.WaitGroup){

	for {
		url:= <-in
		res,err :=http.Get(url)
		if err !=nil{
			fmt.Printf("error %v",err.Error())
		}else{
			fmt.Printf("GET %s: %d\n",url,res.StatusCode)
		}
        wg.Done()
	}
}

func main() {
	var wg sync.WaitGroup
    work:=make(chan string,1024)
    numWorker:=6000

    for i:=0;i<numWorker;i++{

    	go WebGetWorker(work,&wg)
	}
	urls:=[]string {"http://example.com","http://packtpub.com","http://reddit.com","http://twitter.com","http://facebook.com"}
	for i:=0;i<numWorker;i++{
		for _,url :=range urls{
			wg.Add(1)
			work<-url
		}
	}
	wg.Wait()
}
