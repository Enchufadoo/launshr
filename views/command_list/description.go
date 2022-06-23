package command_list

import (
	"fmt"
	"launshr/parser"
)

func RenderDescription(node parser.CommandNode) string {
	resultString := ""

	if node.IsParent() {
		return resultString
	}

	resultString += fmt.Sprintf("%s\n\n", node.Command)

	if node.WorkingDirectory != "" {
		resultString += fmt.Sprintf("Working Directory: %s", node.WorkingDirectory)
	}

	return resultString
}
