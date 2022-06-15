package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type CommandTestCase struct {
	json        string
	node        CommandNode
	description string
}

func TestInvalidJsonFile(t *testing.T) {
	got, err := ParseConfigFile([]byte("ffdsfs{]"))

	assert.Equal(t, CommandNode{}, got)
	assert.NotNil(t, err)
}

func TestMissingCommandThrowsError(t *testing.T) {
	t.Skip("Not implemented, I think I will leave it as the parent's node name, " +
		"rather than assuming an error")
	_, err := ParseConfigFile([]byte(`
		{
			"test" : {
				"name" : "%s"
			}
		}
	`))

	assert.NotNil(t, err)
}

func TestMissingNameUsesCommandAsName(t *testing.T) {
	command := "norm"
	got, err := ParseConfigFile([]byte(`
		{
			"test" : {
				"command" : "norm"
			}
		}
	`))

	node := CommandNode{
		Nodes: []CommandNode{
			{Command: command, Name: command},
		},
	}

	assert.Nil(t, err)

	if !reflect.DeepEqual(node, got) {
		t.Errorf("got %v, want %v", got, node)
	}
}

func TestParseConfigFile(t *testing.T) {
	testCases := []CommandTestCase{
		{"{}", CommandNode{}, "Testing an empty file"},
		singleCommandCase(),
		multipleCommandsCase(),
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {

			got, err := ParseConfigFile([]byte(test.json))

			assert.Nil(t, err)

			if !reflect.DeepEqual(test.node, got) {
				t.Errorf("got %v, want %v", got, test.node)
			}
		})
	}
}

func singleCommandCase() CommandTestCase {
	command := "docker compose rm something"
	name := "Docker RM something"

	json := fmt.Sprintf(`{
		"test" : {
			"command": "%s",
			"name" : "%s"
		}
	}`, command, name)

	node := CommandNode{
		Nodes: []CommandNode{
			{Command: command, Name: name},
		},
	}

	return CommandTestCase{json, node, "Configuration with a single command"}
}

func multipleCommandsCase() CommandTestCase {
	command1 := "cd something"
	name1 := "change dir"

	command2 := "ls"
	name2 := "List something"

	json := fmt.Sprintf(`{
		"a" : {
			"command": "%s",
			"name" : "%s"
		},
		"b" : {
			"command": "%s",
			"name" : "%s"
		}
	}`, command1, name1, command2, name2)

	node := CommandNode{
		Nodes: []CommandNode{
			{Command: command1, Name: name1},
			{Command: command2, Name: name2},
		},
	}

	return CommandTestCase{json, node, "Configuration with multiple commands"}
}
