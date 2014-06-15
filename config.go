package main

import (
	"encoding/xml"	
	"os"
	"io/ioutil"
	"crypto/md5"
	"strings"
	"bytes"
	"log"
)

type ServerConfig struct{
	DirPath string
	XMLName xml.Name `xml:"serverconfig"`
	Courses []*Course `xml:"course"`
	metaSum []byte
}


func LoadServerConfig (dir string)(*ServerConfig,error){

	xmlFile, err := os.Open(dir+"/meta.xml")
	if err != nil {
		return nil,err
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)
	h := md5.New()
	

	config:=new(ServerConfig)

	config.metaSum=h.Sum(b)
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


func (sc *ServerConfig) IsUpdated()(bool){
	
	dir:=sc.DirPath
	newConfig,err:=LoadServerConfig(dir)
	if err!=nil{
		return false
	}
	
	if bytes.Compare(newConfig.metaSum,sc.metaSum)!=0 {
		log.Printf("Detected update in meta.xml file\n")
		return true
	}
	
	if len(newConfig.Courses)!=len(sc.Courses){
		log.Printf("Detected different number of courses\n")
		return true
	}

	for c:=range newConfig.Courses{
		if len(newConfig.Courses[c].Tasks)!=len(sc.Courses[c].Tasks){
			log.Printf("Detected course with different number of tasks\n")
			return true
		}
		for t:= range newConfig.Courses[c].Tasks{
			newHash:=newConfig.Courses[c].Tasks[t].taskSum
			oldHash:=sc.Courses[c].Tasks[t].taskSum
			if bytes.Compare(newHash,oldHash)!=0 {
				log.Printf("Detected update in task '%s'\n",sc.Courses[c].Tasks[t].Id)
				return true
			}
		}
	}

	return false
}






