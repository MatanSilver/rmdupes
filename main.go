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
	LightFileInfoWrapper lightfileinfowrapper
	Hash string
}

type lightfileinfowrapper struct {
	Info os.FileInfo
	Path string
}

func filesInDirectory(dir string) []lightfileinfowrapper {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		if len(files) == 0 {
			return []lightfileinfowrapper{}
		}
		log.Fatalf("Error reading directory: %v\n\t%v", dir, err)
	}
	var fileinfos []lightfileinfowrapper
	for _, file := range files {
		if file.IsDir() {
			//if the file is a directory, recursively add the directory's contents
			log.Printf("Recursing into %s\n", file.Name())
			if file.Name() != ".git" && file.Name() != ".DS_Store" { // ignore .git folder to avoid spam
				newfileinfos := filesInDirectory(strings.Join([]string{dir, file.Name()}, "/"))
				fileinfos = append(fileinfos, newfileinfos...)
			}
		} else {
			fullpath := strings.Join([]string{dir, file.Name()}, "/")
			if file.Size() > 10000000 {
				log.Printf("\tSkipping large file %s\n", fullpath)
				break
			}
			fileinfo := lightfileinfowrapper{file, fullpath}
			fileinfos = append(fileinfos, fileinfo)
		}
	}
	return fileinfos
}

func computeFileHashes(files []lightfileinfowrapper) []fileinfowrapper {
	var fileinfos []fileinfowrapper
	for _, file := range files {
		log.Printf("\tHashing file %s\n", file.Path)
		fileinfo := computeFileHash(file)
		fileinfos = append(fileinfos, fileinfo)
	}
	return fileinfos
}

func computeFileHash(file lightfileinfowrapper) fileinfowrapper {
	buff, err := ioutil.ReadFile(file.Path)
	if err != nil {
		log.Fatal(err)
	}
	hasher := sha256.New()
	hasher.Write(buff)
	fileinfo := fileinfowrapper{file,hex.EncodeToString(hasher.Sum(nil))}
	return fileinfo
}

func rmworker(jobs <-chan string) {
	for j := range jobs {
		err := os.Remove(j)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func RmDupes(dryrun bool, path string, concurrency int) {
	fileinfos := computeFileHashes(filesInDirectory(path))
	fileinfos_map := map[string]string{} //hash to path
	var totalsize int64 = 0
	jobs := make(chan string)
	for i := 1; i <= concurrency; i++ {
		go rmworker(jobs)
	}
	for _, fileinfo := range fileinfos {
		if _, ok := fileinfos_map[fileinfo.Hash]; !ok {
			fileinfos_map[fileinfo.Hash] = fileinfo.LightFileInfoWrapper.Path
		} else {
			log.Printf("File flagged for removal: %s\n\tExisting file: %s\n", fileinfo.LightFileInfoWrapper.Path, fileinfos_map[fileinfo.Hash])
			totalsize += fileinfo.LightFileInfoWrapper.Info.Size()

			if !dryrun {
				jobs <- fileinfo.LightFileInfoWrapper.Path
				//os.Remove(fileinfo.Path)
			}
		}
	}
	close(jobs)
	log.Printf("Total space cleared: %d\n", totalsize)
}

func main() {
	app := cli.NewApp()
	app.Name = "rmdupes"
	app.Version = "0.0.1"
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
		cli.IntFlag{
			Name:  "concurrency, c",
			Usage: "Number of threads to run deleting files",
			Value: 1,
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
		RmDupes(c.Bool("dry-run"), c.String("directory"), c.Int("concurrency"))
		return nil
	}
	app.Run(os.Args)
}
