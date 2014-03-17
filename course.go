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
	Desc string `xml:"description"`
	Tasks []Task
}


type Task struct{
	Name string
	Desc string
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
	}

	return config
}