package parser

const ConfigNodeName = "$config"

type ConfigNode struct {
	WorkingDirectory string
}

func IsConfigNode(key string) bool {
	return key == ConfigNodeName
}
