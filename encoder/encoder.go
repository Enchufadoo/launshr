package encoder

import (
	"github.com/google/uuid"
	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"launshr/config"
	"launshr/parser"
	"launshr/utils"
	"launshr/views/add_node"
	"launshr/views/edit_node"
	"reflect"
	"strings"
)

func EditAndSaveToFile(msg edit_node.EditCommandMsg, filePath string) {

	msg.Node.Name = msg.Msg.Name
	msg.Node.Command = msg.Msg.Command
	msg.Node.WorkingDirectory = msg.Msg.WorkingDirectory

	err := saveEdit(msg.Node, filePath)
	if err != nil {
		utils.ExitError(utils.CouldNotSaveEditJSONFile, err)
	}
}

func AddAndSaveToFile(msg add_node.AddCommandMsg, filePath string) *parser.CommandNode {

	n := parser.CommandNode{
		Command:          msg.Msg.Command,
		Name:             msg.Msg.Name,
		WorkingDirectory: msg.Msg.WorkingDirectory,
		Parent:           msg.Node,
	}

	err := saveAdd(&n, filePath)
	if err != nil {
		utils.ExitError(utils.CouldNotSaveJSONFile, err)
	}

	nodes := config.GetConfig(filePath)
	return &nodes
}

func saveAdd(node *parser.CommandNode, configFilePath string) error {
	joinedKey := ""
	if node.Parent.JsonFullKey != nil {
		joinedKey = strings.Join(node.Parent.JsonFullKey[:], ".") + "."
	}

	id := uuid.New()

	joinedKey += id.String()

	newKeysMap := map[string]string{
		parser.NameKey:             "Name",
		parser.CommandKey:          "Command",
		parser.WorkingDirectoryKey: "WorkingDirectory",
	}
	var err error

	content, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		return err
	}
	stringContent := string(content)

	for k, v := range newKeysMap {
		nodeValue := getCommandNodeField(node, v)
		if nodeValue != "" {
			stringContent, err = sjson.Set(stringContent, joinedKey+"."+k, nodeValue)
			if err != nil {
				return err
			}
		}
	}

	prettyJson := pretty.Pretty([]byte(stringContent))

	err = ioutil.WriteFile(configFilePath, prettyJson, 0644)

	if err != nil {
		return err
	}

	return nil
}

func saveEdit(node *parser.CommandNode, configFilePath string) error {
	joinedKey := strings.Join(node.JsonFullKey[:], ".")

	newKeysMap := map[string]string{
		parser.NameKey:             "Name",
		parser.CommandKey:          "Command",
		parser.WorkingDirectoryKey: "WorkingDirectory",
	}
	var err error

	content, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		return err
	}
	stringContent := string(content)

	for k, v := range newKeysMap {
		nodeValue := getCommandNodeField(node, v)
		if nodeValue != "" {
			stringContent, err = sjson.Set(stringContent, joinedKey+"."+k, nodeValue)
			if err != nil {
				return err
			}
		} else {
			stringContent, err = sjson.Delete(stringContent, joinedKey+"."+k)
			if err != nil {
				return err
			}
		}
	}

	err = ioutil.WriteFile(configFilePath, []byte(stringContent), 0644)

	if err != nil {
		return err
	}

	return nil
}

func getCommandNodeField(v *parser.CommandNode, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}
