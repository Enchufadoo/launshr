package main

import (
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"launshr/config"
	"launshr/utils"
	"launshr/views/main_view"
	"os"
)

func main() {

	configFilePath, err := getConfigFilePath()

	if err != nil {
		utils.ExitError(utils.NoConfigFileProvided, err)
	}

	nodes := config.GetConfig(configFilePath)

	mainView := main_view.New(&nodes, configFilePath)

	p := tea.NewProgram(mainView, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Printf("Error running launshr: %v", err)
		p.Quit()
	}
}

func getConfigFilePath() (string, error) {
	if len(os.Args) < 2 || os.Args[1] == "" {
		return "", errors.New("you need to pass a configuration file as argument")
	}

	configFilePath := os.Args[1]

	return configFilePath, nil
}
