package commandlist

import (
	"github.com/stretchr/testify/assert"
	"launshr/parser"
	"testing"
)

func TestRenderDescription(t *testing.T) {
	cases := []struct {
		description string
		node        parser.CommandNode
		want        string
	}{
		{
			"Getting a description with full node data",
			parser.CommandNode{Name: "Hola", Command: "cd hola", WorkingDirectory: "/"},
			"cd hola\n\nWorking Directory: /",
		},
		{
			"Parent node gets an empty description",
			parser.CommandNode{Name: "Hola", Parent: &parser.CommandNode{}},
			"",
		},
		{
			"No name and working directory just shows the command",
			parser.CommandNode{Command: "cd hi"},
			"cd hi\n\n",
		},
	}

	for _, test := range cases {
		t.Run(test.description, func(t *testing.T) {
			got := RenderDescription(test.node)

			assert.Equal(t, test.want, got)
		})
	}

}
