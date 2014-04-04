package main

import (
	//"fmt"
	"github.com/russross/blackfriday"
	"strings"
	"regexp"
)


func GetContentTitle(content []byte)(string){
	headerReg := regexp.MustCompile("(?m)^\\# .+\\r?\\n")
	out:=string(headerReg.Find(content))
	return strings.Trim(out,"# \t\n")
}


func GetProperty(content []byte, key string)(string){
	propReg := regexp.MustCompile("(?mi)^@>"+key+":.+\\r?\\n")
	p:=string(propReg.Find(content))
	p=strings.ToLower(p)
	f:=strings.Split(p,":")
	if ((f==nil) || (len(f)<2)){
		return ""
	}
	return strings.ToLower(strings.Trim(f[1]," \t\r\n"))
}


func GetContentHTML(content []byte)(string){
	// Clean all properties
	propReg := regexp.MustCompile("(?mi)^@>.+$")
	out:=propReg.ReplaceAll(content,[]byte(""))
	output := blackfriday.MarkdownCommon(out)
	return string(output)
}
