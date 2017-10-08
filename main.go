package main

import (
	"crypto/sha256"
	"encoding/hex"
	//"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type FileInfoWrapper struct {
	Info os.FileInfo
	Path string
	Hash string
}

// Takes in a directory path. Recursively crawls the directory and outputs a
// list of paths of files in that directory and subdirectories
func ls(dir string) []FileInfoWrapper {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		//log.Fatalf("dir: %s, err: %s", dir, err)
		log.Println(err)
	}
	var fileinfos []FileInfoWrapper
	for _, file := range files {
		if file.IsDir() {
			//if the file is a directory, recursively add the directory's contents
			if file.Name() != ".git" { //kindly ignore .git folder to avoid spam
				//log.Printf("entering directory: %s\n", file.Name())
				fileinfos = append(fileinfos, ls(strings.Join([]string{dir, file.Name()}, "/"))...)
			}
		} else {
			fullpath := strings.Join([]string{dir, file.Name()}, "/")
			//now we generate a hash, which might be useful for checking for
			//duplicates
			buff, err := ioutil.ReadFile(fullpath)
			if err != nil {
				log.Fatal(err)
			}
			hasher := sha256.New()
			hasher.Write(buff)
			if err != nil {
				log.Fatal(err)
			}
			fileinfo := FileInfoWrapper{file, fullpath, hex.EncodeToString(hasher.Sum(nil))}
			fileinfos = append(fileinfos, fileinfo)
		}
	}
	return fileinfos
}

func main() {
	fileinfos := ls(".")
	fileinfos_map := map[string]string{}
	for _, fileinfo := range fileinfos {
		if _, ok := fileinfos_map[fileinfo.Hash]; !ok {
			fileinfos_map[fileinfo.Hash] = fileinfo.Path
		} else {
			os.Remove(fileinfo.Path)
		}
	}
}