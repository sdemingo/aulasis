package main

import "log"


func main() {
	srv,err:=CreateServer("./courses")
	if err!=nil{
		log.Panic(err)
	}
	srv.Start()
}

