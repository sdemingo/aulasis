package main


import (
	"text/template"
	"net/http"
	"fmt"
	"strings"
	"os"
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
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("srv/resources")))) 
	http.HandleFunc("/courses/",srv.coursesHandler)
	http.ListenAndServe(":9090", nil)
}






func (srv *Server) coursesHandler(w http.ResponseWriter, r *http.Request) {

	rpath:=strings.TrimPrefix(r.URL.Path,"/courses/")


	if rpath=="" || strings.HasSuffix(rpath,".html"){
		// Aplicar template
		fields:=strings.Split(rpath,"/")
		if (len(fields)==1){
			if fields[0]==""{
				renderTemplate(w,r,"index",srv.Config)
			}else{
				cname:=strings.TrimSuffix(fields[0],".html")
				c:=srv.Config.GetCourseById(cname)
				if c==nil{
					renderTemplate(w,r,"error",nil)
					return
				}
				renderTemplate(w,r,"course",c)
			}
		}else if (len(fields)==2){
			fields[1]=strings.TrimSuffix(fields[1],".html")

			course,task:=fields[0],fields[1]
			c:=srv.Config.GetCourseById(course)
			if c==nil{
				renderTemplate(w,r,"error",nil)
				return
			}
			t:=c.GetTaskById(task)
			if t==nil{
				renderTemplate(w,r,"error",nil)
				return
			}
			renderTemplate(w,r,"task",t)
		}

	}else{
		rpath="srv/courses/"+rpath
		info,err:=os.Stat(rpath)
		fmt.Printf("%s\n",rpath)
		if err!=nil || info.IsDir(){
			renderTemplate(w,r,"error",nil)
			return
		}
		http.ServeFile(w, r, rpath)
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


