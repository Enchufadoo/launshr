package config

import (
	"fmt"
	"io/ioutil"
	"launshr/parser"
	"launshr/utils"
)

type Reader interface {
	ReadFile(path string) ([]byte, error)
}

type FileReader struct{}

func (r FileReader) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func openConfig(configFilePath string, r Reader) ([]byte, error) {

	content, errReadFile := r.ReadFile(configFilePath)

	if errReadFile != nil {
		fmt.Println("Failed to open configuration file: " + configFilePath)
		fmt.Println(errReadFile)
		return []byte{}, errReadFile
	}

	return content, errReadFile
}

func GetConfig(configFilePath string) parser.CommandNode {
	configFileContent, err := openConfig(configFilePath, FileReader{})

	if err != nil {
		utils.ExitError(utils.CouldNotOpenConfigFile, err)
	}

	jsonParser := parser.JsonParser{}
	nodes, err := jsonParser.ParseConfigFile(configFileContent)

	if err != nil {
		utils.ExitError(utils.CouldNotParseJSONFile, "Could not parse JSON file", err)
	}

	return nodes
}
