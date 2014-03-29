package main

import "fmt"

func main() {
	srv,err:=CreateServer("./srv")
	if err!=nil{
		panic(fmt.Sprintf("%v\n",err))
	}
	srv.Start()
}

