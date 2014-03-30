package main

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)



func IterDirectory(rootPath string, dirPath string, z *zip.Writer )(error) {
	dir, err := os.Open(dirPath)
	if err!=nil{
		return err
	}
	defer dir.Close()

	fis, err := dir.Readdir( 0 )
	if err!=nil{
		return err
	}

	for _, fi := range fis {
		curPath := dirPath + "/" + fi.Name()
		if fi.IsDir() {
			err:=IterDirectory(rootPath, curPath, z )
			if err != nil {
				return err
			}
		} else {
			f, err := z.Create(strings.TrimPrefix(curPath,rootPath))
			if err != nil {
				return err
			}

			b, err := ioutil.ReadFile(curPath)
			if err != nil { 
				return err
			}
			_, err = f.Write(b)
			if err != nil {
				return err
			}
			
		}
	}
	return nil
}


// rootPath will be removed from the abs path before adding it to the zip file

func Zip(outFilePath string, rootPath string, inPath string )(error) {

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	z := zip.NewWriter(buf)

	err:=IterDirectory(rootPath, inPath, z )
	if err!=nil{
		return err
	}

	err = z.Close()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outFilePath, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}



/*
func main(){
	targetFilePath := "/tmp/test.zip"
	inputDirPath := "/tmp/test"
	Zip(targetFilePath, inputDirPath, inputDirPath)
}*/