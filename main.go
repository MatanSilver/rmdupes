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
	//fmt.Println("test")
	/*_ = sha256.New()
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang",
			Value: "english",
			Usage: "language for the greeting",
		},
	}

	app.Action = func(c *cli.Context) error {
		return nil
	}

	app.Run(os.Args)
	*/
	fileinfos := ls(".")
	fileinfos_map := map[string]string{}
	for _, fileinfo := range fileinfos {
		if _, ok := fileinfos_map[fileinfo.Hash]; !ok {
			fileinfos_map[fileinfo.Hash] = fileinfo.Path
		} else {
			/*
					      log.Printf("[0]%s is a duplicate of [1]%s\nEnter [0,1] to pick a file to keep, or 2 to keep both", fileinfo.Path, fileinfos_map[fileinfo.Hash])
								var input string
								fmt.Scanln(&input)
								choice, err := strconv.Atoi(input)
								if err != nil {
									log.Println("Error reading input: defaulting to keep both files")
									choice = 2
								}

				choice := 0
				switch choice {
				case 0:
					os.Remove(fileinfo.Path)
				case 1:
					os.Remove(fileinfos_map[fileinfo.Hash])
					fileinfos_map[fileinfo.Hash] = fileinfo.Path
				}
			*/
			os.Remove(fileinfo.Path)
		}
	}
	//fmt.Printf("%v\n", fileinfos_map)
}
