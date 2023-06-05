package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"

	"github.com/yogisinha/filesfinder"
)

func main() {

	//filesfinder.RunCLI()

	// t, err := time.Parse(time.DateTime, "2020-04-26 11:32:04")
	// fmt.Println(t, err)
	// t1, err := time.Parse(time.DateTime, "2020-04-26 13:32:04")
	// fmt.Println(t.Before(t1))
	// fmt.Println(t1.Before(t))

	path := "/home/yogi/github.com/yogisinha/files/files.zip"
	zr, err := zip.OpenReader(path)
	defer zr.Close()
	if err != nil {
		log.Fatal("Error : ", err)
	}
	f, _ := filesfinder.New(
		zr,
		filesfinder.WithDays(10),
	)
	results, err := f.OlderThan()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, res := range results {
		fmt.Printf("%s		%s\n", res.Path, res.ModTime)
	}
	fmt.Println("Erros:")
	for _, res := range results {
		if res.Error != nil {
			fmt.Printf("%s\n", res.Error)
		}
	}

}
