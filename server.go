package main


import (
	"text/template"
	"net/http"
	"strings"
	"os"
	"io"
	"time"
	"regexp"
	"log"
)

const MaxBytesBodySize=20*1024*1024
const CoursesDir="./srv/courses"



type Server struct{
	DirPath string
	Config *ServerConfig
	tmpl *template.Template
}



func CreateServer(dirpath string)(*Server, error){

	srv:=new(Server)
	srv.DirPath=strings.TrimRight(dirpath,"/")

	_,err:=os.Stat(dirpath+"/courses")
	if err!=nil{
		return nil,err
	}

	_,err=os.Stat(dirpath+"/resources")
	if err!=nil{
		return nil,err
	}

	config,err:=LoadServerConfig(srv.DirPath)
	if err!=nil{
		return nil,err
	}

	srv.Config=config
	
	return srv,nil
}


func (srv *Server) Start(){
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("srv/resources")))) 

	http.HandleFunc("/package/",srv.packageHandler)
	http.HandleFunc("/courses/",srv.coursesHandler)
	http.HandleFunc("/submit/",srv.submitHandler)
	http.HandleFunc("/",srv.rootHandler)

	http.ListenAndServe(":9090", nil)
}


func (srv *Server) getCourseAndTask(url string)(*Course,*Task){
	fields:=strings.Split(url,"/")
	var c *Course
	var t *Task

	if (len(fields)==1){
		if fields[0]!=""{
			course:=strings.TrimSuffix(fields[0],".html")
			c=srv.Config.GetCourseById(course)
			return c,nil
		}
	}else if (len(fields)==2){
		fields[1]=strings.TrimSuffix(fields[1],".html")
		course,task:=fields[0],fields[1]
		c=srv.Config.GetCourseById(course)
		if c==nil{
			return nil,nil
		}
		t=c.GetTaskById(task)
		return c,t
	}
	
	return nil,nil
}


func cleanName(name string)(string){
	var noChars = regexp.MustCompile("[^A-Za-záéíóúÁÉÍÓÚñÑüÜ]+")
	out:=noChars.ReplaceAllString(name,"")
	out=strings.ToLower(out)
	return out
}


func (srv *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/courses", http.StatusMovedPermanently)
}



func (srv *Server) submitHandler(w http.ResponseWriter, r *http.Request) {

	rpath:=strings.TrimPrefix(r.URL.Path,"/submit/")

	var task *Task
	if isDinamycUrl(rpath){
		_,task=srv.getCourseAndTask(rpath)
	}else{
		errorHandler(w,r,"Bad submit request")
	}

	name:=cleanName(strings.ToLower(r.FormValue("name")))
	surname:=cleanName(r.FormValue("surname"))

	if name=="" || surname==""{
		errorHandler(w,r,"Bad parametres in request")
		return
	}

	dir:=CoursesDir+"/"+task.Course.Id+"/"+task.Id+"/submitted/"+name+"-"+surname
	err:=os.MkdirAll(dir,0755)	
	if name=="" || surname==""{
		errorHandler(w,r,err.Error())
		return
	}

	//parse the multipart form in the request
	err = r.ParseMultipartForm(MaxBytesBodySize)
	if err != nil {
		errorHandler(w,r,err.Error())
		return
	}
	
	//get a ref to the parsed multipart form
	m := r.MultipartForm
	
	//get the *fileheaders
	files := m.File["files"]
	if err != nil {
		errorHandler(w,r,err.Error())
		return
	}

	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			errorHandler(w,r,err.Error())
			return
		}


		//create destination file making sure the path is writeable.
		dst, err := os.Create(dir+"/"+files[i].Filename)
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}


		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			errorHandler(w,r,err.Error())
			return
		}		
	}

	sr:=new (SubmitReport)
	sr.Files=len(files)
	sr.Stamp=time.Now()
	sr.Addr=getRequestIP(r)
	sr.Task=task

	log.Printf("Task %s submitted from %s\n",task.Title,getRequestIP(r))
	renderTemplate(w,r,"submitted",sr)
}



func (srv *Server) coursesHandler(w http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.Path,"submit"){
		errorHandler(w,r,"Acceso no autorizado")
		return
	}
	rpath:=strings.TrimPrefix(r.URL.Path,"/courses/")
	if rpath==""{
		renderTemplate(w,r,"index",srv.Config)
		return
	}

	if isDinamycUrl(rpath){
		course,task:=srv.getCourseAndTask(rpath)
		if task!=nil{
			renderTemplate(w,r,"task",task)
			return
		}

		if course!=nil{
			renderTemplate(w,r,"course",course)
			return
		}
	
		errorHandler(w,r,"recurso desconocido")
		return
	}else{
		rpath="srv/courses/"+rpath
		info,err:=os.Stat(rpath)
		if err!=nil || info.IsDir(){
			errorHandler(w,r,err.Error())
			return
		}
		http.ServeFile(w, r, rpath)
	}
}


func (srv *Server) packageHandler(w http.ResponseWriter, r *http.Request){

	log.Printf("Packaging request from %s\n",getRequestIP(r))
	rpath:=strings.TrimPrefix(r.URL.Path,"/package/")
	if isDinamycUrl(rpath){
		_,task:=srv.getCourseAndTask(rpath)
		if task!=nil{
			taskfile,err:=task.Package()
			if err!=nil{
				log.Printf("Error: %v\n",err.Error())
			}

			w.Header().Set("Content-Disposition", "attachment; filename="+task.Id+".zip")
			w.Header().Set("Content-type", "application/zip")
			http.ServeFile(w, r, taskfile)

			defer os.Remove(taskfile)

			return
		}
	}
	errorHandler(w,r,"Tarea desconocida para paquetar")
}


func errorHandler(w http.ResponseWriter, r *http.Request, message string){

	log.Printf("Error: %s from %s\n",message,getRequestIP(r))
	t := template.Must(template.ParseFiles("views/error.html"))
	err:=t.Execute(w, message)
	if err!=nil{
		log.Printf("%v\n",err)
	}
}


func renderTemplate(w http.ResponseWriter, r *http.Request, 
	name string,
	cont interface{}) {

	t := template.Must(template.ParseFiles("views/"+name+".html"))
	err:=t.Execute(w, cont)
	if err!=nil{
		log.Printf("%v\n",err)
	}
}



func isDinamycUrl(url string)(bool){
	//ends with an .html extensions or is the root 
	return url=="" || strings.HasSuffix(url,".html")
}


func getRequestIP(r *http.Request)(string){
	f:=strings.Split(r.RemoteAddr,":")
	return f[0]
}