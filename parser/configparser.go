package parser

import (
	"github.com/buger/jsonparser"
)

type CommandNode struct {
	Nodes            []CommandNode
	Parent           *CommandNode
	Name             string
	Command          string
	WorkingDirectory string
}

func (node *CommandNode) IsParent() bool {
	return node.Parent != nil
}

func parseJSONStructure(elementValue []byte, elementType jsonparser.ValueType, node *CommandNode) error {
	switch elementType {
	case jsonparser.Object:

		handler := func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			switch dataType {
			case jsonparser.Object:
				if dataType == jsonparser.Object {
					newNode, err := createNode(value, key, dataType, node)
					if err != nil {
						return err
					}
					node.Nodes = append(node.Nodes, newNode)
				}
			}

			return nil
		}
		err := jsonparser.ObjectEach(elementValue, handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func ParseConfigFile(configFileContent []byte) (CommandNode, error) {
	var commandTree CommandNode
	value, elementType, _, err := jsonparser.Get(configFileContent)

	if err != nil {
		return commandTree, err
	}

	err = parseJSONStructure(value, elementType, &commandTree)
	if err != nil {
		return commandTree, err
	}

	return commandTree, nil
}

func createNode(value []byte, key []byte, dataType jsonparser.ValueType,
	node *CommandNode) (CommandNode, error) {
	newNode := CommandNode{}

	name, nameErr := jsonparser.GetString(value, "name")
	command, commandErr := jsonparser.GetString(value, "command")
	workingDirectory, workingDirectoryErr := jsonparser.GetString(value, "wd")

	if workingDirectoryErr != nil {
		newNode.WorkingDirectory = workingDirectory
	}

	if commandErr == nil {
		if nameErr != nil {
			name = command
		}

		newNode.Name = name
		newNode.Command = command
	} else {
		if name == "" {
			name = string(key)
		}

		newNode.Name = "\u2630 " + name
		newNode.Parent = node
		err := parseJSONStructure(value, dataType, &newNode)

		if err != nil {
			return newNode, err
		}
	}

	return newNode, nil
}
