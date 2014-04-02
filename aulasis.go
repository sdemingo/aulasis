package main

import "log"


func main() {
	srv,err:=CreateServer("./srv")
	if err!=nil{
		log.Panic(err)
	}
	srv.Start()
}

