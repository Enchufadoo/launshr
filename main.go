package main

import (
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"io/ioutil"
	"launshr/parser"
	"launshr/views/main_view"
	"os"
)

type Reader interface {
	ReadFile(path string) ([]byte, error)
}

type FileReader struct{}

func (r FileReader) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func main() {
	configFileContent, configFilePath, err := openConfig(FileReader{})

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

	mainView := main_view.New(&nodes, configFilePath)

	mainView.SetCurrentView(main_view.CommandListView)
	p := tea.NewProgram(mainView, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		p.Quit()
	}
}

func openConfig(r Reader) ([]byte, string, error) {

	if len(os.Args) < 2 || os.Args[1] == "" {
		return []byte{}, "", errors.New("you need to pass a configuration file as argument")
	}

	configFilePath := os.Args[1]

	content, errReadFile := r.ReadFile(configFilePath)

	if errReadFile != nil {
		fmt.Println("Failed to open configuration file: " + configFilePath)
		fmt.Println(errReadFile)
		return []byte{}, configFilePath, errReadFile
	}

	return content, configFilePath, errReadFile
}
