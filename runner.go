package main

import (
	"os/exec"
	"github.com/wsxiaoys/terminal"
	"time"
	"fmt"
	"strings"
	"sync"
)

type Runner struct{
	project *Project
}

func (r *Runner) Start() {
	var wg sync.WaitGroup
	total := len(r.project.global) + len(r.project.local)
	pending, complete := make(chan Process), make(chan string)
	for i:=0;i<total;i++{
		wg.Add(1)
		go r.Run(pending, complete, &wg)
	}
	for _, process := range r.project.global {
		pending <- process
	}
	for _, process := range r.project.local {
		pending <- process
	}

	for message := range complete {
		fmt.Println( message)
	}
	wg.Wait()
}

func (r *Runner) Run(in <-chan Process, out chan<- string, wg *sync.WaitGroup) {
	for process := range in{
		var commands []string
		commands = strings.Split(process.command, " ")

		var args []string
		program := commands[0]
		if(len(commands) > 1){
			args = append(commands[:0], commands[1:]...)
		}else{
			args = []string{}
		}

		cmd := exec.Command(program, args...)

		if process.dir != "" {
			cmd.Dir = process.dir
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println(err)
		}

		err = cmd.Start()
		if err != nil {
			fmt.Println(err)
		}
		ch := make(chan string)
		quit := make(chan bool)
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := stdout.Read(buf)
				if n != 0 {
					ch <- string(buf[:n])
				}
				if err != nil {
					break
				}
			}
			close(ch)
		}()

		time.AfterFunc(time.Second, func() { quit <- true })

		loop:
		for {
			select {
				case lines, ok := <-ch:
					if !ok {
						break loop
					}
					for _, line := range strings.Split(lines, "\n") {
						if(line == ""){
							continue
						}
						terminal.Stderr.
						Color(process.color).Print(fmt.Sprintf("[%s] ", process.name)).
						Reset().Print(line).Nl()
					}
			}
		}

		if err != nil{
			out <- fmt.Sprintf("%s failed to start (%s)", process.name, err)
			wg.Done()
			return
		}
		out <- fmt.Sprintf("%s started. Command: %s. Args: %s", process.name, program, args)
		err = cmd.Wait()
		if err != nil{
			out <- fmt.Sprintf("%s failed to wait (%s)", process.name, err)
			wg.Done()
			return
		}
		out <- fmt.Sprintf("%s finished", process.name)
		wg.Done()
	}
}

