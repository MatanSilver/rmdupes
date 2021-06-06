package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestLSDupes(t *testing.T) {
	os.Mkdir("tmp", 0777)
	os.Mkdir("tmp/subdir", 0777)
	err := ioutil.WriteFile("tmp/file1", []byte("this is a test"), 0777)
	if err != nil {
		log.Fatal(err)
	}
	/*err = ioutil.WriteFile("./tmp/file2", []byte("this is a test"), 0666)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("./tmp/subdir/file1", []byte("this is a test"), 0666)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("./tmp/subdir/file2", []byte("this is another test"), 0666)
	if err != nil {
		log.Fatal(err)
	}*/

	//fileinfos := ls("./tmp")
	//log.Println(fileinfos)

	//teardown
	err = os.RemoveAll("./tmp")
	if err != nil {
		log.Fatal(err)
	}
}
