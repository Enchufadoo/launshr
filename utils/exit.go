package utils

import (
	"fmt"
	"os"
)

const (
	NoConfigFileProvided = iota
	CouldNotOpenConfigFile
	CouldNotParseJSONFile
)

func ExitError(errorCode int, message ...any) {
	fmt.Println(message)
	os.Exit(errorCode)
}
