package parser

import (
	"github.com/google/uuid"
	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"reflect"
	"strings"
)

func SaveAddToFile(node *CommandNode, configFilePath string) error {
	joinedKey := ""
	if node.Parent.JsonFullKey != nil {
		joinedKey = strings.Join(node.Parent.JsonFullKey[:], ".") + "."
	}

	id := uuid.New()

	joinedKey += id.String()

	newKeysMap := map[string]string{
		NameKey:             "Name",
		CommandKey:          "Command",
		WorkingDirectoryKey: "WorkingDirectory",
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

func SaveEditToFile(node *CommandNode, configFilePath string) error {
	joinedKey := strings.Join(node.JsonFullKey[:], ".")

	newKeysMap := map[string]string{
		NameKey:             "Name",
		CommandKey:          "Command",
		WorkingDirectoryKey: "WorkingDirectory",
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

func getCommandNodeField(v *CommandNode, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}
