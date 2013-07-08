package main

import (
	"encoding/json"
	"os/user"
	"io/ioutil"
	"fmt"
	"os"
)

type Settings struct {
	Projects map[string]string `json:"projects"`
}

func LoadSettings() (settings *Settings, err error) {
	homedir := HomeDir()
	content, err := ioutil.ReadFile(homedir + "/.hack/config")

	// Lazily create hack config dir and retry
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
		err = nil
		settings = &Settings{make(map[string]string)}
	}else{
		err = json.Unmarshal(content, &settings)
		if err != nil {
			return
		}
	}

	return
}

func (s *Settings) Write() (err error) {
	homedir := HomeDir()

	// Convert to json
	b, err := json.Marshal(s)
	if err != nil {
		return
	}

	// Create the config dir (no-op)
	err = os.MkdirAll(homedir + "/.hack", 0700)
	if err != nil {
		return
	}

	// Write to the file
    err = ioutil.WriteFile(homedir + "/.hack/config", b, 0700)
	if err != nil {
		return
	}
	return
}

func HomeDir() (string) {
	usr, err := user.Current()
    if err != nil {
        fmt.Println( err )
    }
    return usr.HomeDir
}
