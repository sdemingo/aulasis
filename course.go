package main


import (
	"os"
	"fmt"
	"encoding/xml"
	"io/ioutil"
	"strings"
	"time"
	"text/template"
)



type ServerConfig struct{
	DirPath string
	XMLName xml.Name `xml:"serverconfig"`
	Courses []*Course `xml:"course"`
}


func LoadServerConfig (dir string)(*ServerConfig){

	xmlFile, err := os.Open(dir+"/courses/meta.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)

	config:=new(ServerConfig)

	config.DirPath=dir
	err = xml.Unmarshal(b, &config)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}

	for c:=range config.Courses{
		config.Courses[c].Desc=strings.Trim(config.Courses[c].Desc, " \n")
		LoadCourse(dir, config.Courses[c])
	}

	return config
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

	dirpath:=basedir+"/courses/"+course.Id
	infos,err:=ioutil.ReadDir(dirpath)
	if err != nil {
		fmt.Printf("error: %v", err)
		return 
	}

	course.Tasks=make([]*Task,len(infos))
	course.BaseDir=basedir

	for i:=range infos{
		t:=LoadTask(course,infos[i].Name())
		course.Tasks[i]=t  //maybe nill
	}
}


func (c *Course) GetTaskById(id string)(*Task){
	return LoadTask(c,id)
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


func LoadTask(course *Course,taskId string)(*Task){
	orgfile:=course.BaseDir+"/courses/"+course.Id+"/"+taskId+"/info.org"
	orgFile, err := os.Open(orgfile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer orgFile.Close()

	b, _ := ioutil.ReadAll(orgFile)

	task:=new(Task)
	task.Content=Org2HTML(b,taskId)
	task.Id=taskId
	task.Title=ParseHeader(b,"TITLE")
	task.Course=course
	task.Status=ParseProperty(b,"status")
	if task.Status==""{
		task.Status=TASK_CLOSED_STATUS
	}

	return task
}


func (task *Task) CheckStatus(st string)(bool){
	return task.Status==st
}



func (task *Task) Package()(file string,err error){

	file=""
	out,err:=ioutil.TempFile("",task.Id+"-pack")
	if err!=nil{
		return
	}
	tdir,err:=ioutil.TempDir("",task.Id)
	if err!=nil{
		return
	}
	err=os.Mkdir(tdir+"/"+task.Id,0755)
	if err!=nil{
		return
	}

	dir,err:=os.Open(task.Course.BaseDir+"/courses/"+task.Course.Id+"/"+task.Id)
	if err!=nil{
		return
	}

	names,err:=dir.Readdirnames(-1)
	if err!=nil{
		return
	}

	taskdir:=task.Course.BaseDir+"/courses/"+task.Course.Id+"/"+task.Id
	fmt.Printf("Packaging task %s\n",taskdir)


	f,err:=os.Create(tdir+"/info.html")
	if err!=nil{
		return
	}
	t := template.Must(template.ParseFiles("views/local-task.html"))
	err=t.Execute(f, task)
	if err!=nil{
		fmt.Printf("%v\n",err)
	}
	fmt.Printf("\tAdding %s\n",tdir+"/info.html")


	copyFile(task.Course.BaseDir+"/resources/default.css",tdir+"/"+task.Id+"/default.css")
	fmt.Printf("\tAdding %s\n",tdir+"/"+task.Id+"/default.css")

	/*
	 The Task directories should not have subdirectories. They are
	 ignored during the packaging process
	 */

	for i:=range names{
		file:=taskdir+"/"+names[i]
		info,err:=os.Stat(file)
		if err==nil && info.IsDir()==false{
			fmt.Printf("\tAdding %s\n",file)
			copyFile(file,tdir+"/"+task.Id+"/"+info.Name())
		}else{
			fmt.Printf("\tIgnoring %s\n",file)
		}
	}
	


	// Compress Temp Dir
	file=out.Name()+".zip"
	Zip(file, tdir, tdir)

	// Delete Temp Dir
	err=os.RemoveAll(tdir)
	if err!=nil{
		fmt.Printf("%v\n",err)
	}

	fmt.Printf("Package %s is ready\n",file)
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