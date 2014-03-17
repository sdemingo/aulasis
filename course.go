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


type Course struct{
	Name string `xml:"name"`
	Path string `xml:"path"`
	Desc string `xml:"description"`
	Tasks []*Task
}


type Task struct{
	Name string
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
		fmt.Printf("Loading course %s\n",config.Courses[c].Name)
		LoadCourse(config.Courses[c])
	}

	return config
}



func LoadCourse(course *Course){
	
	dirpath:="./srv/"+course.Path
	infos,err:=ioutil.ReadDir(dirpath)
	if err != nil {
		fmt.Printf("error: %v", err)
		return 
	}

	course.Tasks=make([]*Task,len(infos))

	for i:=range infos{
		path:=dirpath+"/"+infos[i].Name()+"/info.org"
		t:=LoadTask(path)
		course.Tasks[i]=t  //maybe nill
		if (t!=nil){
			fmt.Printf("\tLoading task \"%s\"\n",t.Name)
		}
	}
}



func LoadTask(orgfile string)(*Task){
	orgFile, err := os.Open(orgfile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer orgFile.Close()

	b, _ := ioutil.ReadAll(orgFile)

	task:=new(Task)
	task.Content=b
	task.Name=ParseHeader(b,"TITLE")

	return task
}