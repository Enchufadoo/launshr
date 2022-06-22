package main

import (
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"io/ioutil"
	"launshr/parser"
	"launshr/views/commandlist"
	"os"
)

type Reader interface {
	ReadFile(filename string) ([]byte, error)
}

type FileReader struct{}

func (r FileReader) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func main() {
	configFileContent, err := openConfig(FileReader{})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	jsonParser := parser.JsonParser{}
	nodes, err := jsonParser.ParseConfigFile(configFileContent)

	if err != nil {
		fmt.Println("Could not parse JSON file")
		fmt.Println(err)
		os.Exit(1)
	}

	p := tea.NewProgram(commandlist.InitialModel(&nodes))

	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		p.Quit()
	}
}

func openConfig(r Reader) ([]byte, error) {

	if len(os.Args) < 2 || os.Args[1] == "" {
		return []byte{}, errors.New("you need to pass a configuration file as argument")
	}

	configFile := os.Args[1]

	content, errReadFile := r.ReadFile(configFile)

	if errReadFile != nil {
		fmt.Println("Failed to open configuration file: " + configFile)
		fmt.Println(errReadFile)
		return []byte{}, errReadFile
	}

	return content, errReadFile
}
