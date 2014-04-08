package main


import (
	"os"
	"log"
	"encoding/xml"
	"io/ioutil"
	"strings"
	"time"
	"text/template"
	"path/filepath"
	"regexp"
)



type ServerConfig struct{
	DirPath string
	XMLName xml.Name `xml:"serverconfig"`
	Courses []*Course `xml:"course"`
}


func LoadServerConfig (dir string)(*ServerConfig,error){

	xmlFile, err := os.Open(dir+"/meta.xml")
	if err != nil {
		return nil,err
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)

	config:=new(ServerConfig)

	config.DirPath=dir
	err = xml.Unmarshal(b, &config)
	if err != nil {
		return nil,err
	}

	for c:=range config.Courses{
		config.Courses[c].Desc=strings.Trim(config.Courses[c].Desc, " \n")
		LoadCourse(dir, config.Courses[c])
	}

	return config,nil
}


func (sc *ServerConfig) GetCourseById(id string)(*Course){
	for i:=range sc.Courses{
		if sc.Courses[i].Id==id{
			return sc.Courses[i]
		}
	}
	return nil
}










type Course struct{
	BaseDir string
	Name string `xml:"name"`
	Id string `xml:"path"`
	Desc string `xml:"description"`
	Tasks []*Task
}


func LoadCourse(basedir string, course *Course){

	dirpath:=basedir+"/"+course.Id
	infos,err:=ioutil.ReadDir(dirpath)
	if err != nil {
		log.Printf("Error loading course %s: %v",course.Id, err)
		return 
	}

	course.Tasks=make([]*Task,len(infos))
	course.BaseDir=basedir

	for i:=range infos{
		t,err:=LoadTask(course,infos[i].Name())
		if err!=nil{
			log.Printf("Error loading task %s\n",infos[i].Name())
		}
		course.Tasks[i]=t  //maybe nill
	}
}


func (c *Course) GetTaskById(id string)(*Task){
	t,err:=LoadTask(c,id)
	if err!=nil{
		log.Printf("Error loading task %s\n",id)
	}
	return t
}





const TASK_OPEN_STATUS = "open"
const TASK_CLOSED_STATUS = "close"
const TASK_HIDE_STATUS = "hide"

const TASK_PROP_YES = "yes"
const TASK_PROP_NO = "no"

type Task struct{
	Course *Course
	Title string
	Id string
	Content string
	Status string
	LogFile string
}



type SubmitReport struct{
	Task *Task
	Path string
	Name string
	Surname string
	Addr string
	Stamp time.Time
	Files int
}


func LoadTask(course *Course,taskId string)(*Task,error){
	orgfile:=course.BaseDir+"/"+course.Id+"/"+taskId+"/info.md"
	orgFile, err := os.Open(orgfile)
	if err != nil {
		return nil,err
	}
	defer orgFile.Close()

	b, _ := ioutil.ReadAll(orgFile)

	task:=new(Task)
	task.Content=GetContentHTML(b)
	task.Id=taskId
	task.Title=GetContentTitle(b)
	task.Course=course
	task.LogFile=GetProperty(b,"logfile")
	if task.LogFile!=TASK_PROP_NO &&
		task.LogFile != TASK_PROP_YES{
		task.LogFile=TASK_PROP_NO   //by default
		log.Printf("Error: task file %s has a bad logfile property\n",task.Id)
	}
	task.Status=GetProperty(b,"status")
	if task.Status!=TASK_CLOSED_STATUS && 
		task.Status!=TASK_OPEN_STATUS && 
		task.Status!=TASK_HIDE_STATUS {
		
		task.Status=TASK_CLOSED_STATUS  //by default
		log.Printf("Error: task file %s has a bad status property\n",task.Id)
	}

	return task,nil
}


func (task *Task) CheckStatus(st string)(bool){
	return task.Status==st
}

func (task *Task) WriteLog(msg string){
	if task.LogFile==TASK_PROP_NO{
		return
	}
	
	logfile:=task.Course.BaseDir+"/"+task.Course.Id+"/"+task.Id+"/submits.log"
	f, err := os.OpenFile(logfile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("task log error: %v", err)
	}
	defer f.Close()
}



func (task *Task) Package()(file string,err error){

	file=""

	tdir,err:=ioutil.TempDir("",task.Id)
	if err!=nil{
		return
	}
	err=os.Mkdir(tdir+"/"+task.Id,0755)
	if err!=nil{
		return
	}

	dir,err:=os.Open(task.Course.BaseDir+"/"+task.Course.Id+"/"+task.Id)
	if err!=nil{
		return
	}

	names,err:=dir.Readdirnames(-1)
	if err!=nil{
		return
	}

	taskdir:=task.Course.BaseDir+"/"+task.Course.Id+"/"+task.Id
	f,err:=os.Create(tdir+"/info.html")
	if err!=nil{
		return
	}

	// Replace remote links to local links before apply the template
	linksReg:=regexp.MustCompile("=\"/courses/[^/]+/[^/]+/")
	task.Content=linksReg.ReplaceAllString(task.Content,"=\""+task.Id+"/")

	t := template.Must(template.ParseFiles(ResourcesDir+"/templates/local-task.html"))
	err=t.Execute(f, task)
	if err!=nil{
		return
	}
	copyFile(ResourcesDir+"/css/default.css",tdir+"/"+task.Id+"/default.css")

	/*
	 The Task directories should not have subdirectories. They are
	 ignored during the packaging process
	 */

	for i:=range names{
		file:=taskdir+"/"+names[i]
		info,err:=os.Stat(file)
		if err==nil && info.IsDir()==false{
			copyFile(file,tdir+"/"+task.Id+"/"+info.Name())
		}
	}
	
	file=filepath.Base(tdir)+".zip"
	Zip(file, tdir, tdir)

	defer os.RemoveAll(tdir)

	return				    
}




func copyFile(src,dest string)(error){
	b,err:=ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	err=ioutil.WriteFile(dest,b,0644)
	if err != nil {
		return err
	}

	return nil
}