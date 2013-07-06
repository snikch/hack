package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"path"
	"errors"
	"github.com/wsxiaoys/terminal"
)
func main () {
	var command string
	if(len(os.Args) > 1){
		command = os.Args[1]
	}
	project := &Project{}

	switch command {

		// Run the toil processes
		case "run":
		default:
			err := loadProject("toil.json", project)
			if err != nil {
				PrintError(err)
				return
			}

			dir, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			name := path.Base(dir)
			project.name = name

			runner := &Runner{project: project}
			runner.Start()

		// Generate an empty toilfile
		case "init":
			_, err := ioutil.ReadFile("toil.json")
			if err == nil {
				PrintError(errors.New("Toilfile exists"))
				return
			}
			err = project.Write()
			if err != nil {
				PrintError(err)
				return
			}
			PrintSuccess("Initialized toilfile")

		// Add a process
		case "add":
			if len(os.Args) < 4 {
				PrintError(errors.New("Usage: toil add [name] [command]"))
				return
			}
			err := loadProject("toil.json", project)
			if err != nil {
				PrintError(err)
				return
			}

			name, command := os.Args[2], os.Args[3]
			process := Process{ name: name, command: command}

			global, replaced := IsGlobalCommand(), false
			target := project.Processes(global)

			_, replaced = target[name]
			target[name] = process
			err = project.Write()
			if err != nil {
				panic(err)
			}
			if replaced {
				PrintSuccess(fmt.Sprintf("Replaced %s", name))
			}else{
				PrintSuccess(fmt.Sprintf("Added %s", name))
			}
			return

		// Remove a process
		case "rm":
			if len(os.Args) < 3 {
				PrintError(errors.New("Usage: toil rm [name]"))
				return
			}
			err := loadProject("toil.json", project)
			if err != nil {
				PrintError(err)
				return
			}
			name, global := os.Args[2], IsGlobalCommand()

			target := project.Processes(global)

			_, exists := target[name]
			if !exists {
				PrintError(errors.New(fmt.Sprintf("No process named %s", name)))
				return
			}

			delete(target, name)
			err = project.Write()
			if err != nil {
				panic(err)
			}

			PrintSuccess(fmt.Sprintf("Removed %s", name))
			return

		// List all processes
		case "list":
			err := loadProject("toil.json", project)
			if err != nil {
				PrintError(err)
				return
			}
			fmt.Println(fmt.Sprintf("%d processes in toilfile", project.Size()))
			if len(project.local) > 0 {
				fmt.Println("Local:")
				for _, process := range project.local {
					fmt.Println(fmt.Sprintf(" - %s: %s", process.name, process.command))
				}
			}
			if len(project.global) > 0 {
				fmt.Println("Global:")
				for _, process := range project.global {
					fmt.Println(fmt.Sprintf(" - %s: %s", process.name, process.command))
				}
			}
			return

		// Register the current toilfile globally
		case "register":
			panic(errors.New("TODO: Register not yet implemented"))
			return
	}
}

func PrintlnColored(str string, color string) {
	terminal.Stderr.Color(color).Print(str).Reset().Nl()
	return
}

func PrintSuccess(str string) {
	PrintlnColored(fmt.Sprintf("✔ %s", str), "g")
	return
}

func PrintError(err error) {
	PrintlnColored(fmt.Sprintf("✘ %s", err.Error()), "r")
	return
}

func IsGlobalCommand() (bool) {
	for _, t := range os.Args {
		if t == "-g" {
			 return true
		 }
	}
	return false
}
