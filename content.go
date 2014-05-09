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

	/*
	 // basic rendering without any customization

	 output := blackfriday.MarkdownCommon(out)
	 return string(output)
	 */

	// seting custon flags and extensions from blackfriday/markdown.go

	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	htmlFlags |= blackfriday.HTML_SANITIZE_OUTPUT
	//htmlFlags |= blackfriday.HTML_TOC
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	extensions |= blackfriday.EXTENSION_HEADER_IDS

	output:=blackfriday.Markdown(out, renderer, extensions)
	return string(output)
}
