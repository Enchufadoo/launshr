package parser

type CommandNode struct {
	Nodes            []CommandNode
	Parent           *CommandNode
	Name             string
	Command          string
	WorkingDirectory string
	Config           *ConfigNode
}

func (node *CommandNode) IsParent() bool {
	return node.Parent != nil
}
