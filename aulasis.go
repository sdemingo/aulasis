package main

import "log"


func main() {
	srv,err:=CreateServer("./resources","./courses")
	if err!=nil{
		log.Panic(err)
	}
	srv.Start()
}

