package main


import (
	"text/template"
	"net/http"
	"fmt"
	"strings"
)


type Server struct{
	Config ServerConfig
	tmpl *template.Template
}

func CreateServer(config *ServerConfig)(*Server){
	srv:=new(Server)

	srv.Config=*config

	return srv
}

func (srv *Server) Start(){
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web")))) 
	http.HandleFunc("/courses/", srv.coursesHandler)
	http.ListenAndServe(":9090", nil)
}






func (srv *Server) coursesHandler(w http.ResponseWriter, r *http.Request) {


	coursesPath:="/courses/"
	rpath:=r.URL.Path[len(coursesPath):]

	fields:=strings.Split(rpath,"/")
	if (len(fields)==1){
		if fields[0]==""{
			renderTemplate(w,r,"index",srv.Config)
		}else{
			c:=srv.Config.GetCourseByPath(fields[0])
			if c==nil{
				renderTemplate(w,r,"error",nil)
				return
			}
			renderTemplate(w,r,"course",c)
		}
	}

	if (len(fields)==2){
		// load activity
		fmt.Printf("must load task %s\n",fields[1])
	}
}



func renderTemplate(w http.ResponseWriter, r *http.Request, 
	name string, 
	cont interface{}) {

	t := template.Must(template.ParseFiles("views/"+name+".html"))
	err:=t.Execute(w, cont)
	if err!=nil{
		fmt.Printf("%v\n",err)
	}
}


