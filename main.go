package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

type fileinfowrapper struct {
	Info os.FileInfo
	Path string
	Hash string
}

// Takes in a directory path. Recursively crawls the directory and outputs a
// list of paths of files in that directory and subdirectories
func ls(dir string) []fileinfowrapper {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		if len(files) == 0 {
			return []fileinfowrapper{}
		}
		log.Fatalf("Error reading directory: %v\n\t%v", dir, err)
	}
	var fileinfos []fileinfowrapper
	for _, file := range files {
		if file.IsDir() {
			//if the file is a directory, recursively add the directory's contents
			log.Printf("Recursing into %s\n", file.Name())
			if file.Name() != ".git" && file.Name() != ".DS_Store" { //kindly ignore .git folder to avoid spam
				//log.Printf("entering directory: %s\n", file.Name())
				newfileinfos := ls(strings.Join([]string{dir, file.Name()}, "/"))
				fileinfos = append(fileinfos, newfileinfos...)
			}
		} else {
			fullpath := strings.Join([]string{dir, file.Name()}, "/")
			if file.Size() > 10000000 {
				log.Printf("\tSkipping large file %s\n", fullpath)
				break
			}

			//now we generate a hash, which might be useful for checking for
			//duplicates
			log.Printf("\tHashing file %s\n", fullpath)
			buff, err := ioutil.ReadFile(fullpath)
			if err != nil {
				log.Fatal(err)
			}
			hasher := sha256.New()
			hasher.Write(buff)
			//TODO possibly remove use of hex and just put to []byte
			fileinfo := fileinfowrapper{file, fullpath, hex.EncodeToString(hasher.Sum(nil))}
			fileinfos = append(fileinfos, fileinfo)
		}
	}
	return fileinfos
}

func RmDupes(dryrun bool, path string) {
	fileinfos := ls(path)
	fileinfos_map := map[string]string{} //hash to path
	var totalsize int64 = 0
	for _, fileinfo := range fileinfos {
		if _, ok := fileinfos_map[fileinfo.Hash]; !ok {
			fileinfos_map[fileinfo.Hash] = fileinfo.Path
		} else {
			log.Printf("File flagged for removal: %s\n\tExisting file: %s\n", fileinfo.Path, fileinfos_map[fileinfo.Hash])
			totalsize += fileinfo.Info.Size()
			if !dryrun {
				os.Remove(fileinfo.Path)
			}
		}
	}
	log.Printf("Total space cleared: %d\n", totalsize)
}

func main() {
	app := cli.NewApp()
	app.Name = "rmdupes"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Only print logs of files to be deleted or errors",
		},
		cli.StringFlag{
			Name:  "directory, d",
			Usage: "Specify parent directory in which to remove duplicates",
			Value: ".",
		},
		/*
				cli.BoolFlag{
					Name:	"recurse|r",
					Usage:	"Recursively search for duplicates inside subdirectories",
				},
			}
		*/
	}
	app.Usage = "Removes duplicate files by _content_"
	app.Action = func(c *cli.Context) error {
		RmDupes(c.Bool("dry-run"), c.String("directory"))
		return nil
	}
	app.Run(os.Args)
}
