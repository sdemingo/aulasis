package main


import (
	"os"
	"fmt"
	"encoding/xml"
	"io/ioutil"
	"strings"
)


type ServerConfig struct{
	XMLName xml.Name `xml:"serverconfig"`
	Courses []*Course `xml:"course"`
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
	Name string `xml:"name"`
	Id string `xml:"path"`
	Desc string `xml:"description"`
	Tasks []*Task
}


type Task struct{
	Name string
	Id string
	Content []byte
}




func LoadServerConfig (metafile string)(*ServerConfig){

	xmlFile, err := os.Open(metafile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)


	config:=new(ServerConfig)

	err = xml.Unmarshal(b, &config)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}

	for c:=range config.Courses{
		config.Courses[c].Desc=strings.Trim(config.Courses[c].Desc, " \n")
		//fmt.Printf("Loading course %s\n",config.Courses[c].Name)
		LoadCourse(config.Courses[c])
	}

	return config
}



func LoadCourse(course *Course){
	
	dirpath:="./srv/"+course.Id
	infos,err:=ioutil.ReadDir(dirpath)
	if err != nil {
		fmt.Printf("error: %v", err)
		return 
	}

	course.Tasks=make([]*Task,len(infos))

	for i:=range infos{
		t:=LoadTask(dirpath,infos[i].Name())
		course.Tasks[i]=t  //maybe nill
	}
}



func LoadTask(coursedir,taskname string)(*Task){
	orgfile:=coursedir+"/"+taskname+"/info.org"
	orgFile, err := os.Open(orgfile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer orgFile.Close()

	b, _ := ioutil.ReadAll(orgFile)

	task:=new(Task)
	task.Content=b
	task.Id=taskname
	task.Name=ParseHeader(b,"TITLE")

	return task
}