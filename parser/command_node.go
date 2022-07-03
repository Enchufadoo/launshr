package parser

const (
	NameKey             = "name"
	CommandKey          = "command"
	WorkingDirectoryKey = "wd"
)

type CommandNode struct {
	Nodes            []CommandNode
	Parent           *CommandNode
	Name             string
	Command          string
	WorkingDirectory string
	Config           *ConfigNode
	JsonFullKey      []string
}

func (node *CommandNode) IsParent() bool {
	return node.Parent != nil
}
