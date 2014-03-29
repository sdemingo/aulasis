package main

import (
	"archive/zip"
	"bytes"
	"log"
	"io/ioutil"
	"os"
	"fmt"
	"strings"
)



func IterDirectory(rootPath string, dirPath string, z *zip.Writer ) {
	dir, err := os.Open(dirPath)
	if err!=nil{
		log.Fatal(err)
	}
	defer dir.Close()
	fis, err := dir.Readdir( 0 )
	if err!=nil{
		log.Fatal(err)
	}

	for _, fi := range fis {
		curPath := dirPath + "/" + fi.Name()
		if fi.IsDir() {
			IterDirectory(rootPath, curPath, z )
		} else {
			fmt.Printf( "adding... %s\n", curPath )

			f, err := z.Create(strings.TrimPrefix(curPath,rootPath))
			if err != nil {
				log.Fatal(err)
			}

			b, err := ioutil.ReadFile(curPath)
			if err != nil { panic(err) }


			_, err = f.Write(b)
			if err != nil {
				log.Fatal(err)
			}
			
		}
	}
}


// rootPath will be removed from the abs path before adding it to the zip file

func Zip(outFilePath string, rootPath string, inPath string ) {

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	z := zip.NewWriter(buf)

	IterDirectory(rootPath, inPath, z )

	err := z.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(outFilePath, buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println( "zip ok" )
}



/*
func main(){
	targetFilePath := "/tmp/test.zip"
	inputDirPath := "/tmp/test"
	Zip(targetFilePath, inputDirPath, inputDirPath)
}*/