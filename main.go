package main

import (
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
	configFile := os.Args[1]

	configFileContent, err := openConfig(FileReader{}, configFile)

	if err != nil {
		fmt.Println("Failed opening "+configFile, err)
		os.Exit(1)
	}

	nodes, err := parser.ParseConfigFile(configFileContent)

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

func openConfig(r Reader, filename string) ([]byte, error) {
	content, errReadFile := r.ReadFile("./" + filename)
	return content, errReadFile
}
