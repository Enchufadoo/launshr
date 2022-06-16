package commandlist

import (
	"fmt"
	"launshr/parser"
)

func RenderDescription(node parser.CommandNode) string {
	resultString := ""

	if node.IsParent() {
		return resultString
	}

	if node.Name != "" {
		resultString += fmt.Sprintf("%s \n", node.Name)
	}

	resultString += fmt.Sprintf("%s \n", node.Command)

	if node.WorkingDirectory != "" {
		resultString += fmt.Sprintf("Working Directory: %s \n", node.Name)
	}

	return resultString
}
