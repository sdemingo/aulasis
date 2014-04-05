package main

import (
	"log"
	"bitbucket.org/kardianos/osext"
	"strings"
	"flag"
)

var portFlag = flag.Int("p", 9090, "Service port")
var docFlag = flag.String("d", "", "Courses and tasks directory")

func main() {

	flag.Parse()
	port:=*portFlag

	execpath,err:=osext.ExecutableFolder()
	if err!=nil{
		log.Panic(err)
	}	
	execpath=strings.TrimRight(execpath,"/\\")

	docpath:=strings.TrimRight(*docFlag,"/\\")
	if docpath==""{
		docpath=execpath+"/courses"
	}

	srv,err:=CreateServer(execpath+"/resources",docpath)
	if err!=nil{
		log.Panic(err)
	}
	srv.Start(port)
}

