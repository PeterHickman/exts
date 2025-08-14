package main

import (
	"flag"
	"fmt"
	ep "github.com/PeterHickman/expand_path"
	"github.com/PeterHickman/toolbox"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var root string
var search_hidden bool
var extensions map[string]int

func init() {
	var a = flag.Bool("a", false, "Search inside hidden directories")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatalln("No directory supplied")
	}

	search_hidden = *a

	root = flag.Arg(0)

	root, _ = ep.ExpandPath(root)
	if !toolbox.FileExists(root) {
		fmt.Println("The directory " + flag.Arg(0) + " does not exist")
		os.Exit(2)
	}

	extensions = make(map[string]int)
}

func extract_ext(filename string) {
	base := filepath.Base(filename)
	parts := strings.Split(base, `.`)

	if len(parts) == 2 && parts[0] == "" {
		// Skip things like .gitignore
		return
	}

	if len(parts) > 1 {
		e := parts[len(parts)-1]
		extensions[e]++
	}
}

func hidden_directory(path string) bool {
	base := filepath.Base(path)
	return strings.HasPrefix(base, ".")
}

func main() {
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			i, _ := os.Stat(path)

			if i.IsDir() {
				if hidden_directory(path) && !search_hidden {
					return filepath.SkipDir
				}
			} else {
				extract_ext(path)
			}

			return nil
		})

	if err != nil {
		fmt.Println(err)
		os.Exit(6)
	}

	fmt.Println("extension,count")
	for k, v := range extensions {
		fmt.Printf("%s,%d\n", k, v)
	}
}
