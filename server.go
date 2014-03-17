package main


import (
	"text/template"
	"net/http"
	"fmt"
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
	course:=r.URL.Path[len(coursesPath):]

	if course!=""{
		c:=srv.Config.GetCourseByPath(course)
		if c==nil{
			renderTemplate(w,r,"error",nil)
			return
		}

		renderTemplate(w,r,"course",c)
	}else{
		renderTemplate(w,r,"index",srv.Config)
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


