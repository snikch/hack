package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"path"
)
func main () {

	p := os.Getenv("PATH")

	fmt.Println(fmt.Sprintf("Path: %s", p))
	content, err := ioutil.ReadFile("toil.json")

	if err != nil {
		fmt.Println(fmt.Sprintf("Could not %s", err))
		return
	}

	project := &Project{}
	jsonToProject(content, project)

	dir, err := os.Getwd()
	if err != nil{
		panic(err)
	}

	name := path.Base(dir)
	project.name = name

	runner := &Runner{project: project}

	runner.Start()
	fmt.Println(project, runner)
}
