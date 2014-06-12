package main


import (
	"text/template"
	"net/http"
	"net"
	"strings"
	"os"
	"io"
	"time"
	"regexp"
	"log"
	"fmt"
	"math/rand"
	"sync"
)

const MaxBytesBodySize=20*1024*1024
var ResourcesDir string

const ErrorTaskNotSubmitted="Tarea no entregada"
const ErrorResourceUnknown="Pedido recurso desconocido"
const ErrorServerSuspended="Servidor parado temporalmente. Inténtelo pasados unos segundos"


type Server struct{
	ResourcesPath string
	DirPath string
	Config *ServerConfig
	tmpl *template.Template
	sWorker chan string
	suspended bool
	reqGroup sync.WaitGroup
}



func CreateServer(respath string,dirpath string)(*Server, error){

	srv:=new(Server)
	srv.DirPath=strings.TrimRight(dirpath,"/")
	srv.ResourcesPath=respath
	ResourcesDir=respath

	_,err:=os.Stat(dirpath)
	if err!=nil{
		return nil,err
	}

	_,err=os.Stat(srv.ResourcesPath)
	if err!=nil{
		return nil,err
	}

	config,err:=LoadServerConfig(srv.DirPath)
	if err!=nil{
		return nil,err
	}

	srv.Config=config
	srv.sWorker=make(chan string)
	srv.suspended=false

	return srv,nil
}


func (srv *Server) Start(port int){
	
	srv.showPublicIp(port)

	go srv.submitWorker()
	go srv.updateWorker()

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir(srv.ResourcesPath)))) 

	http.HandleFunc("/package/",srv.packageHandler)
	http.HandleFunc("/courses/",srv.coursesHandler)
	http.HandleFunc("/submit/",srv.submitHandler)
	http.HandleFunc("/",srv.rootHandler)

	log.Printf("**** Aulasis running on port %d ****\n",port)
	http.ListenAndServe(fmt.Sprintf(":%d",port), nil)
}

func (srv *Server) showPublicIp(port int){
	addrs,err:=net.InterfaceAddrs()
	if err!=nil{
		return
	}
	fmt.Printf("  \n");
	fmt.Printf("  Open your browse and try to connect to following addresses on port %d:\n\n",port)
	for _,a:=range addrs{
		fmt.Printf("\t%s\n",a.String())
	}
	fmt.Printf("\n");
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


func (srv *Server) renderTemplate(w http.ResponseWriter, r *http.Request, 
	name string,
	cont interface{}) {

	t := template.Must(template.ParseFiles(srv.ResourcesPath+"/templates/"+name+".html"))
	err:=t.Execute(w, cont)
	if err!=nil{
		log.Printf("render error: %v\n",err)
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


// Go routine to order the submit requests on the filesystem
// and check if the submit directory exits
func (srv *Server) submitWorker(){

	for ;;{
		subpath:=<-srv.sWorker
		_,err:=os.Stat(subpath)
		if err!=nil{
			srv.sWorker<-subpath
		}else{
			rand.Seed(time.Now().UTC().UnixNano())
			s:=fmt.Sprintf("%d",rand.Int())
			srv.sWorker<-subpath+"-"+s
		}
	}
}

// Goroutine to check configserver updates.
func (srv *Server) updateWorker(){
	for ;;{
		time.Sleep(30 * time.Second)
		if srv.Config.IsUpdated(){
			srv.suspended=true
			// if IsUpdate return true it must manage the update
			// loading the new config when all reques are done
			srv.reqGroup.Wait()
			log.Printf("todas las peticiones están terminadas")
			log.Printf("Ejecutar ahora la nueva carga");
		}
	}
}



func (srv *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/courses", http.StatusMovedPermanently)
}



/*
 Handlers
*/

func (srv *Server) submitHandler(w http.ResponseWriter, r *http.Request) {

	rpath:=strings.TrimPrefix(r.URL.Path,"/submit/")

	if srv.suspended{
		srv.errorHandler(w,r,ErrorServerSuspended,nil)
		return
	}

	srv.reqGroup.Add(1)
	defer srv.reqGroup.Done()

	var task *Task
	if isDinamycUrl(rpath){
		_,task=srv.getCourseAndTask(rpath)
	}else{
		srv.errorHandler(w,r,ErrorTaskNotSubmitted,nil)
		return
	}

	name:=cleanName(strings.ToLower(r.FormValue("uname")))
	surname:=cleanName(r.FormValue("surname"))

	if name=="" || surname==""{
		srv.errorHandler(w,r,ErrorTaskNotSubmitted,nil)
		return
	}

	dir:=srv.DirPath+"/"+task.Course.Id+"/"+task.Id+"/submitted/"+name+"-"+surname
	srv.sWorker<-dir
	dir=<-srv.sWorker // waiting for submitWorker confirmation
	err:=os.MkdirAll(dir,0755)	
	if err!=nil{
		srv.errorHandler(w,r,ErrorTaskNotSubmitted,err)
		return
	}

	err = r.ParseMultipartForm(MaxBytesBodySize)
	if err != nil {
		srv.errorHandler(w,r,ErrorTaskNotSubmitted,err)
		return
	}
	m := r.MultipartForm

	files := m.File["files"]
	if err != nil {
		srv.errorHandler(w,r,ErrorTaskNotSubmitted,err)
		return
	}

	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			srv.errorHandler(w,r,ErrorTaskNotSubmitted,err)
			return
		}

		dst, err := os.Create(dir+"/"+files[i].Filename)
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(dst, file); err != nil {
			srv.errorHandler(w,r,ErrorTaskNotSubmitted,err)
			return
		}		
	}

	sr:=new (SubmitReport)
	sr.Files=len(files)
	sr.Stamp=time.Now()
	sr.Addr=getRequestIP(r)
	sr.Task=task

	msg:=fmt.Sprintf("%s %s entrega %d ficheros desde %s",
		r.FormValue("name"), r.FormValue("surname"), sr.Files, sr.Addr)
	task.WriteLog(msg)

	log.Printf("Task '%s' submitted from %s\n",task.Title,getRequestIP(r))
	srv.renderTemplate(w,r,"submitted",sr)
}



func (srv *Server) coursesHandler(w http.ResponseWriter, r *http.Request) {

	if srv.suspended{
		srv.errorHandler(w,r,ErrorServerSuspended,nil)
		return
	}

	srv.reqGroup.Add(1)
	defer srv.reqGroup.Done()

	if strings.Contains(r.URL.Path,"submit"){
		srv.errorHandler(w,r, ErrorResourceUnknown,nil)
		return
	}

	rpath:=strings.TrimPrefix(r.URL.Path,"/courses/")
	if rpath==""{
		srv.renderTemplate(w,r,"index",srv.Config)
		return
	}

	if isDinamycUrl(rpath){
		course,task:=srv.getCourseAndTask(rpath)
		if task!=nil{
			srv.renderTemplate(w,r,"task",task)
			return
		}

		if course!=nil{
			srv.renderTemplate(w,r,"course",course)
			return
		}
		
		srv.errorHandler(w,r, ErrorResourceUnknown,nil)
		return
	}else{
		rpath=srv.DirPath+"/"+rpath
		info,err:=os.Stat(rpath)
		if err!=nil || info.IsDir(){
			srv.errorHandler(w,r, ErrorResourceUnknown,err)
			return
		}
		http.ServeFile(w, r, rpath)
	}
}


func (srv *Server) packageHandler(w http.ResponseWriter, r *http.Request){

	if srv.suspended{
		srv.errorHandler(w,r,ErrorServerSuspended,nil)
		return
	}

	srv.reqGroup.Add(1)
	defer srv.reqGroup.Done()

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
	srv.errorHandler(w,r, ErrorResourceUnknown, nil)
}


func (srv *Server) errorHandler(w http.ResponseWriter, r *http.Request, message string,erri error){

	if erri!=nil{
		log.Printf("Error: %s from %s\n",erri,getRequestIP(r))
	}else{
		log.Printf("Error: %s from %s\n",message,getRequestIP(r))
	}

	t := template.Must(template.ParseFiles(srv.ResourcesPath+"/templates/error.html"))
	err:=t.Execute(w, message)
	if err!=nil{
		log.Printf("render error: %v\n",err)
	}
}
