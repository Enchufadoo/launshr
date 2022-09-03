package parser

import (
	"github.com/buger/jsonparser"
)

type JsonParser struct{}

func (c *JsonParser) parseStructure(elementValue []byte, elementType jsonparser.ValueType,
	node *CommandNode, fullKey []string) error {

	config := c.GetConfigNode(elementValue)

	if config != (ConfigNode{}) {
		node.Config = &config
	}

	switch elementType {
	case jsonparser.Object:
		handler := func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			if !IsConfigNode(string(key)) {
				switch dataType {
				case jsonparser.Object:
					if dataType == jsonparser.Object {
						newNode, err := c.createNode(value, key, dataType, node, &config, fullKey)
						if err != nil {
							return err
						}
						node.Nodes = append(node.Nodes, newNode)
					}
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

func (c *JsonParser) ParseConfigFile(configFileContent []byte) (CommandNode, error) {
	var commandTree CommandNode
	value, elementType, _, err := jsonparser.Get(configFileContent)

	if err != nil {
		return commandTree, err
	}

	var fullKey []string
	err = c.parseStructure(value, elementType, &commandTree, fullKey)
	if err != nil {
		return commandTree, err
	}

	return commandTree, nil
}

func (c *JsonParser) createNode(value []byte, key []byte, dataType jsonparser.ValueType,
	node *CommandNode, config *ConfigNode, fullKey []string) (CommandNode, error) {
	newNode := CommandNode{}

	fullKey = append(fullKey, string(key))
	newNode.JsonFullKey = fullKey

	name, _ := jsonparser.GetString(value, NameKey)
	command, _ := jsonparser.GetString(value, CommandKey)

	if config != nil && *config != (ConfigNode{}) {
		newNode.Config = config
	} else if node.Config != nil && *node.Config != (ConfigNode{}) {
		newNode.Config = node.Config
	}

	newNode.WorkingDirectory = c.getWorkingDirectory(&newNode, value)

	if command != "" {
		if name == "" {
			name = command
		}
		newNode.Name = name
		newNode.Command = command
	} else {
		if name == "" {
			name = string(key)
		}

		newNode.Name = name
		newNode.Parent = node
		err := c.parseStructure(value, dataType, &newNode, fullKey)

		if err != nil {
			return newNode, err
		}
	}

	return newNode, nil
}

// If there is a parent node with a working directory, use that
// else use the key wd to get the working directory
func (c *JsonParser) getWorkingDirectory(node *CommandNode, value []byte) string {
	workingDirectory, _ := jsonparser.GetString(value, WorkingDirectoryKey)
	if workingDirectory == "" {
		if node.Config != nil {
			workingDirectory = node.Config.WorkingDirectory
		}
	}

	return workingDirectory
}

// GetConfigNode Parse the $config object for the children nodes to inherit
func (c *JsonParser) GetConfigNode(value []byte) ConfigNode {
	configNode := ConfigNode{}
	value, dataType, _, err := jsonparser.Get(value, ConfigNodeName)

	if err == nil {
		if dataType == jsonparser.Object {
			wd, _ := jsonparser.GetString(value, WorkingDirectoryKey)
			title, _ := jsonparser.GetString(value, ConfigTitleKey)

			configNode.WorkingDirectory = wd
			configNode.Title = title
		}
	}

	return configNode
}
