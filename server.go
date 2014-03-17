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
	http.HandleFunc("/courses/", srv.coursesHandler)
	http.ListenAndServe(":9090", nil)
}






func (srv *Server) coursesHandler(w http.ResponseWriter, r *http.Request) {


	coursesPath:="/courses/"
	course:=r.URL.Path[len(coursesPath):]

	if course!=""{
		c:=srv.Config.GetCourseByPath(course)
		if c==nil{
			t := template.Must(template.ParseFiles("views/error.html"))
			err:=t.Execute(w, nil)
			if err!=nil{
				fmt.Printf("%v\n",err)
			}
		}else{

			t := template.Must(template.ParseFiles("views/course.html"))
			err:=t.Execute(w, c)
			if err!=nil{
				fmt.Printf("%v\n",err)
			}
		}
		
	}else{
		t := template.Must(template.ParseFiles("views/index.html"))
		err:=t.Execute(w, srv.Config)
		if err!=nil{
			fmt.Printf("%v\n",err)
		}
	}
}




