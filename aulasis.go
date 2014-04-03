package main

import (
	"log"
	"bitbucket.org/kardianos/osext"
	"strings"
)

func main() {

	execpath,err:=osext.ExecutableFolder()
	if err!=nil{
		log.Panic(err)
	}	
	execpath=strings.TrimRight(execpath,"/\\")

	// By now, courses are in the same folder that exec file but
	// in the future it has been input by a flag
	docpath:=execpath+"/courses"

	srv,err:=CreateServer(execpath+"/resources",docpath)
	if err!=nil{
		log.Panic(err)
	}
	srv.Start()
}

