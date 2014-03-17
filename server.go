package main


import (
	"html/template"
	"net/http"
	"fmt"
)


type Server struct{
	Config ServerConfig
}

func CreateServer(config *ServerConfig)(*Server){
	srv:=new(Server)

	srv.Config=*config
	return srv
}

func (srv *Server) Start(){
	http.HandleFunc("/courses/", srv.coursesHandler)
	http.ListenAndServe(":9090", nil)
}








func (srv *Server) coursesHandler(w http.ResponseWriter, r *http.Request) {


	coursesPath:="/courses/"
	course:=r.URL.Path[len(coursesPath):]

	if course!=""{
		fmt.Printf("cargamos curso %s\n",course)
		
	}else{
		t, err:= template.ParseFiles("views/index.html")
		
		if err!=nil {
			fmt.Printf("%v",err)
			return
		}

		t.Execute(w, srv.Config)
	}
}




