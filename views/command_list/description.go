package command_list

import (
	"fmt"
	"launshr/parser"
	"os"
)

func RenderDescription(node parser.CommandNode) string {
	resultString := ""

	if node.IsParent() {
		return resultString
	}

	resultString += fmt.Sprintf("%s\n\n", node.Command)

	wdString := ""

	if node.WorkingDirectory != "" {
		wdString = node.WorkingDirectory
	} else {
		cwd, err := os.Getwd()
		if err == nil {
			wdString += cwd + "\n"
		}
		wdString += "Current Folder"
	}

	resultString += fmt.Sprintf("Working Directory \n%s", wdString)

	return resultString
}
