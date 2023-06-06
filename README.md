# files
Small package to find files older than a specific time

### Example program using this package
```
package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"

	"github.com/yogisinha/filesfinder"
)

func main() {

	path := "/path/to/some/folder"
	fs := os.DirFS(path)
	f, _ := filesfinder.New(
		fs, // first argument is fs.FS type

		/* second is the list of time options letting the user specify the time
		as a combination of years, months, days, hours and minutes */
		filesfinder.WithDays(10),
		filesfinder.WithHours(5),
	)
	// above example will find files older than 10 days, 5 hrs.
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

	// Also zip file path can be used.
	path = "/path/to/some/zipfile.zip"
	zr, err := zip.OpenReader(path)
	defer zr.Close()
	if err != nil {
		log.Fatal("Error : ", err)
	}
	f, _ = filesfinder.New(
		zr,
		filesfinder.WithDays(10),
	)
	results, err = f.OlderThan()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, res := range results {
		fmt.Printf("%s		%s\n", res.Path, res.ModTime)
	}

}

```
