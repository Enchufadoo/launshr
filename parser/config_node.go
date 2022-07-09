package parser

const ConfigNodeName = "$config"
const ConfigTitleKey = "title"

type ConfigNode struct {
	WorkingDirectory string
	Title            string
}

func IsConfigNode(key string) bool {
	return key == ConfigNodeName
}
