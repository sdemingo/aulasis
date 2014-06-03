package main

import (
	"log"
	"bitbucket.org/kardianos/osext"
	"strings"
	"flag"
	"fmt"
)

var LICENSE=`
  Copyright (C) 2014  Sergio de Mingo <sdemingo@gmail.com>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.`

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

	fmt.Println(LICENSE)

	srv,err:=CreateServer(execpath+"/resources",docpath)
	if err!=nil{
		log.Panic(err)
	}
	srv.Start(port)
}

