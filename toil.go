package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"path"
)
func main () {
	command := os.Args[1]
	project := &Project{}

	switch command {
		default:
			content, err := ioutil.ReadFile("toil.json")

			if err != nil {
				fmt.Println(fmt.Sprintf("Could not %s", err))
				return
			}

			jsonToProject(content, project)

			dir, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			name := path.Base(dir)
			project.name = name

			runner := &Runner{project: project}
			runner.Start()
		case "init":
			project.Write()
		case "add"
	}
	//fmt.Println(project, runner)
}
