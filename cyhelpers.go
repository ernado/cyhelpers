package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var names = [...]string{
	"cydev",
	"ernado",
	".",
}

var gopath string

func completion(prefixes []string) string {
	candidates := []string{}
	exists := make(map[string]bool)
	for _, prefix := range prefixes {
		for _, name := range names {
			dir := filepath.Join(prefix, name)
			dirs, err := ioutil.ReadDir(dir)
			if err != nil {
				continue
			}
			// do not complete all gopath
			if (name == ".") && strings.HasPrefix(prefix, gopath) {
				continue
			}
			for _, d := range dirs {
				if !d.IsDir() {
					continue
				}
				candidate := d.Name()
				if exists[candidate] {
					candidate = filepath.Join(name, candidate)
				}
				exists[candidate] = true
				candidates = append(candidates, candidate)
			}
		}
	}
	return strings.Join(candidates, " ")
}

func main() {
	var gopathExists bool
	gopath, gopathExists = os.LookupEnv("GOPATH")
	if !gopathExists {
		gopath = "/go"
	}
	prefixes := []string{
		filepath.Join(gopath, "src", "github.com"), // golang projects
		"/src", // other projects
	}
	if len(os.Args) < 2 {
		// returning completion list
		fmt.Print(completion(prefixes))
		os.Exit(1)
	}
	project := os.Args[1] // cydev/project, ernado/project or just project
	candidates := []string{
		project,
		filepath.Join("cydev", project),
		filepath.Join("ernado", project),
	}
	var path string
	for _, prefix := range prefixes {
		for _, candidate := range candidates {
			pathCandidate := filepath.Join(prefix, candidate)
			if _, err := os.Stat(pathCandidate); err == nil {
				path = pathCandidate
				break
			}
		}
	}

	if len(path) == 0 {
		// trying to return CWD
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(dir)
		fmt.Fprintf(os.Stderr, "error: project %s not found\n", project)
		os.Exit(-1)
	}
	fmt.Print(path)
}
