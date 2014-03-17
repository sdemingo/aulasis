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
	Courses []Course `xml:"course"`
}


type Course struct{
	Name string `xml:"name"`
	Path string `xml:"path"`
	Desc string `xml:"description"`
	Tasks []Task
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

		fmt.Printf("Loading course %s from %s/\n",config.Courses[c].Name,config.Courses[c].Path)

		LoadCourse(config.Courses[c].Path)
	}

	return config
}



func LoadCourse(dir string)(*Course){
	
	dirpath:="./srv/"+dir
	infos,err:=ioutil.ReadDir(dirpath)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}

	
	for i:=range infos{
		path:=dirpath+"/"+infos[i].Name()+"/info.org"
		LoadTask(path)
	}


	return nil
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

	fmt.Printf("\tLoading task titled %s from %s\n",task.Name,orgfile)

	return task
}