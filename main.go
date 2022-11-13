package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RepositoryXML struct {
	XMLName      xml.Name        `xml:"repository"`
	Dependencies []DependencyXML `xml:"dependency"`
}

type DependencyXML struct {
	XMLName  xml.Name `xml:"dependency"`
	Name     string   `xml:"name,attr"`
	Url      string   `xml:"url,attr"`
	Revision string   `xml:"revision,attr"`
}

type Dependency struct {
	Name     string
	Url      string
	Revision string
}

func readRepositoryXML(dir string) *RepositoryXML {
	xmlFile, err := os.Open(filepath.Join(dir, "repository.xml"))
	if err != nil {
		return nil
	}
	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		panic(err)
	}

	var repository RepositoryXML
	err = xml.Unmarshal(byteValue, &repository)
	if err != nil {
		panic(err)
	}

	return &repository
}

func pathToDeps() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, "deps")
}

func checkDepsDirectory(depsDir string) {
	if stat, err := os.Stat(depsDir); os.IsNotExist(err) {
		if err = os.Mkdir(depsDir, 0750); err != nil {
			panic(err)
		}
	} else if !stat.IsDir() {
		panic(fmt.Sprintf("%v should be directory", depsDir))
	}
}

func joinStrings(s ...string) string {
	r := ""
	for _, n := range s {
		r = fmt.Sprintf("%s %s", r, n)
	}
	return r
}

type OutputType int

const (
	CaptureOutput  OutputType = 1
	RedirectOutput            = 2
)

func executeGit(depsDir string, outType OutputType, args ...string) (string, error) {
	command := exec.Command("git", args...)
	command.Dir = depsDir
	outString := ""

	if outType == RedirectOutput {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin
		err := command.Run()
		if err != nil {
			// fmt.Printf("run command 'git%s' error: %v", joinStrings(args...), err)
			return "", nil
		}
	} else {
		out, err := command.CombinedOutput()
		outString = string(out)
		if err != nil {
			// fmt.Printf("run command 'git%s' error: %v\noutput: %s", joinStrings(args...), err, outString)
			return outString, err
		}
		fmt.Printf("process 'git%s' output: %s\n", joinStrings(args...), outString)
	}

	return outString, nil
}

func alreadyProcessed(current DependencyXML, list []DependencyXML) bool {
	for i := 0; i < len(list); i++ {
		if list[i].Name != current.Name {
			continue
		}

		if list[i].Revision == current.Revision && list[i].Url == current.Url {
			return true
		}

		panic(fmt.Sprintf("dependency missmatch: %v and %v", current, list[i]))
	}

	return false
}

func main() {
	depsProcessed := []DependencyXML{}
	directories := []string{"./"}
	var currentDirectory string

	depsDir := pathToDeps()
	checkDepsDirectory(depsDir)

	for len(directories) != 0 {
		currentDirectory, directories = directories[0], directories[1:]

		repository := readRepositoryXML(currentDirectory)
		if repository == nil {
			fmt.Printf("-----------------\n")
			fmt.Printf("no repository.xml in %v\n", currentDirectory)
			continue
		}

		for i := 0; i < len(repository.Dependencies); i++ {
			fmt.Printf("-----------------\n")
			dep := repository.Dependencies[i]

			if len(dep.Revision) == 0 {
				dep.Revision = "master"
			}

			if alreadyProcessed(dep, depsProcessed) {
				continue
			}

			depDir := filepath.Join(depsDir, dep.Name)

			fmt.Printf("%v: %v#%v\n", depDir, dep.Url, dep.Revision)

			out, err := executeGit(depsDir, CaptureOutput, "clone", dep.Url, dep.Name)
			if err != nil {
				if !strings.Contains(out, "already exists and is not an empty directory") {
					panic(out)
				}
			}

			_, err = executeGit(depDir, RedirectOutput, "checkout", dep.Revision)
			if err != nil {
				panic(err)
			}

			_, err = executeGit(depDir, RedirectOutput, "pull")
			if err != nil {
				panic(err)
			}

			directories = append(directories, depDir)
			depsProcessed = append(depsProcessed, dep)
		}
	}

	fmt.Printf("-----------------\n")
	fmt.Printf("all dependencies fetched!\n")
}
