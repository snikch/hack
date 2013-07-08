package main

import (
	"github.com/wsxiaoys/terminal"
	"io/ioutil"
	"errors"
	"fmt"
	"os"
	"path"
)

type Cli struct {
	Command string
	Project *Project
}

func NewCli() *Cli {
	command := "run"
	if(len(os.Args) > 1){
		command = os.Args[1]
	}

	return &Cli{ command, &Project{} }
}

func (c *Cli) run() {
	err := c.beforeFilter()
	if err != nil {
		return;
	}
	switch c.Command {
		// Run the toil processes
		case "run":
			dir, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			name := path.Base(dir)
			c.Project.name = name

			runner := &Runner{project: c.Project}
			runner.Start()

		// Generate an empty toilfile
		case "init":
			_, err := ioutil.ReadFile("toil.json")
			if err == nil {
				c.PrintError(errors.New("Toilfile exists"))
				return
			}
			err = c.Project.Write()
			if err != nil {
				c.PrintError(err)
				return
			}
			c.PrintSuccess("Initialized toilfile")

		// Add a process
		case "add":
			if len(os.Args) < 4 {
				c.PrintError(errors.New("Usage: toil add [name] [command]"))
				return
			}

			name, command := os.Args[2], os.Args[3]
			process := Process{ name: name, command: command}

			global, replaced := c.IsGlobalCommand(), false
			target := c.Project.Processes(global)

			_, replaced = target[name]
			target[name] = process
			err = c.Project.Write()
			if err != nil {
				panic(err)
			}
			if replaced {
				c.PrintSuccess(fmt.Sprintf("Replaced %s", name))
			}else{
				c.PrintSuccess(fmt.Sprintf("Added %s", name))
			}
			return

		// Remove a process
		case "rm":
			if len(os.Args) < 3 {
				c.PrintError(errors.New("Usage: toil rm [name]"))
				return
			}

			name, global := os.Args[2], c.IsGlobalCommand()

			target := c.Project.Processes(global)

			_, exists := target[name]
			if !exists {
				c.PrintError(errors.New(fmt.Sprintf("No process named %s", name)))
				return
			}

			delete(target, name)
			err = c.Project.Write()
			if err != nil {
				panic(err)
			}

			c.PrintSuccess(fmt.Sprintf("Removed %s", name))
			return

		// List all processes
		case "list":
			fmt.Println(fmt.Sprintf("%d processes in toilfile", c.Project.Size()))
			if len(c.Project.local) > 0 {
				fmt.Println("Local:")
				for _, process := range c.Project.local {
					fmt.Println(fmt.Sprintf(" - %s: %s", process.name, process.command))
				}
			}
			if len(c.Project.global) > 0 {
				fmt.Println("Global:")
				for _, process := range c.Project.global {
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

func (c *Cli) beforeFilter() (err error) {
	switch c.Command {
		case "run", "add", "rm", "list":
			err = loadProject("toil.json", c.Project)
			if err != nil {
				c.PrintError(err)
			}
			return
	}
	return
}

func (c *Cli) PrintlnColored(str string, color string) {
	terminal.Stderr.Color(color).Print(str).Reset().Nl()
	return
}

func (c *Cli) PrintSuccess(str string) {
	c.PrintlnColored(fmt.Sprintf("✔ %s", str), "g")
	return
}

func (c *Cli) PrintError(err error) {
	c.PrintlnColored(fmt.Sprintf("✘ %s", err.Error()), "r")
	return
}

func (c *Cli) IsGlobalCommand() (bool) {
	for _, t := range os.Args {
		if t == "-g" {
			 return true
		 }
	}
	return false
}
