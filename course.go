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
const TASK_CLOSED_STATUS = "closed"
const TASK_HIDE_STATUS = "hide"


type Task struct{
	Course *Course
	Title string
	Id string
	Content string
	Status string
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
	task.Id=taskId
	task.Title=GetContentTitle(b)
	task.Course=course
	task.Status=GetProperty(b,"status")
	if task.Status==""{
		task.Status=TASK_CLOSED_STATUS
	}
	task.Content=GetContentHTML(b)

	return task,nil
}


func (task *Task) CheckStatus(st string)(bool){
	return task.Status==st
}



func (task *Task) Package()(file string,err error){

	file=""
	/*out,err:=ioutil.TempFile("",task.Id+"-pack")
	if err!=nil{
		return
	}*/
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
	//log.Printf("Packaging task %s\n",taskdir)

	f,err:=os.Create(tdir+"/info.html")
	if err!=nil{
		return
	}
	t := template.Must(template.ParseFiles(ResourcesDir+"/templates/local-task.html"))
	err=t.Execute(f, task)
	if err!=nil{
		return
	}
	//log.Printf("\tAdding %s\n",tdir+"/info.html")

	copyFile(ResourcesDir+"/css/default.css",tdir+"/"+task.Id+"/default.css")
	//log.Printf("\tAdding %s\n",tdir+"/"+task.Id+"/default.css")

	/*
	 The Task directories should not have subdirectories. They are
	 ignored during the packaging process
	 */

	for i:=range names{
		file:=taskdir+"/"+names[i]
		info,err:=os.Stat(file)
		if err==nil && info.IsDir()==false{
			//log.Printf("\tAdding %s\n",file)
			copyFile(file,tdir+"/"+task.Id+"/"+info.Name())
		}else{
			//log.Printf("\tIgnoring %s\n",file)
		}
	}
	
	// Compress Temp Dir
	file=filepath.Base(tdir)+".zip"
	Zip(file, tdir, tdir)

	// Delete Temp Dir
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