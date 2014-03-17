package main

import "fmt"

func main() {

	config:=LoadServerConfig("srv/meta.xml")
	if config==nil{
		fmt.Printf("ServerConfig not loaded\n")
	}


	srv:=CreateServer(config)
	srv.Start()
}

