package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Process struct{
	name string
	command string
	color string
}

type Project struct{
	name string
	global map[string]Process
	local map[string]Process
}

type ExportableProject struct{
	Name string `json:"name"`
	Global map[string]string `json:"global"`
	Local map[string]string `json:"local"`
}

// Write the project to the json file
func (p *Project) Write() (err error) {
	var f *os.File

	f, err = os.OpenFile("toil.json", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0666)
	if err != nil {
		return
    }

    defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(p.AsJson())
	if err != nil {
		return
	}

    return
}

func (p *Project) AsJson() (e *ExportableProject){
	e = &ExportableProject{}
	e.Name = p.name

	e.Local = make(map[string]string)
	for _, process := range p.local {
		e.Local[process.name] = process.command
	}

	e.Global = make(map[string]string)
	for _, process := range p.global {
		e.Global[process.name] = process.command
	}
	return
}

func (p *Project) Size() (i int) {
	i = len(p.global) + len(p.local)
	return
}

func (p *Project) Processes(global bool) (target map[string]Process) {
	if global {
		target = p.global
	}else{
		target = p.local
	}
	return
}
func loadProject(path string, p *Project) (err error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return
	}

	jsonToProject(content, p)
	return
}

func jsonToProject(content []byte, project *Project) {
	var projMap  map[string]*json.RawMessage
	err := json.Unmarshal(content, &projMap)

	var local map[string]string
	if projMap["local"] == nil {
		local = make(map[string]string)
	} else {
		err = json.Unmarshal(*projMap["local"], &local)
		if err != nil {
			return
		}
	}

	var global map[string]string
	if projMap["global"] == nil {
		global = make(map[string]string)
	} else {
		err = json.Unmarshal(*projMap["global"], &global)
		if err != nil {
			return
		}
	}

	colors := []string{"g", "y", "b", "m", "c"}
	numColors := len(colors)
	count := 0

	project.local = make(map[string]Process, len(local))
	for name, command := range local {
		project.local[name] = Process{name, command, colors[count]}
		count++
		if(count > numColors){
			count = 0
		}
	}

	project.global = make(map[string]Process, len(global))
	for name, command := range global {
		project.global[name] = Process{name, command, colors[count]}
		if(count > numColors){
			count = 0
		}
	}
}
