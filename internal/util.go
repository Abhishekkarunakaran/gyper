package internal

import (
	"fmt"
	"regexp"
)

// creates a string slice from the given path string
// eg: for "/v1/private/test" returns []string{"/v1","/private","/test"}
func GetPathList(str string) []string {
	var list []string
	var outString string
	if len(str) == 1 {
		list = append(list, str)
		return list
	}
	for i, ch := range str {
		if ch == '/' {
			list = append(list, outString)
			outString = ""
			outString += string(ch)
			continue
		}
		if i == len(str)-1 {
			outString += string(ch)
			list = append(list, outString)
			continue
		}
		outString += string(ch)

	}
	return list[1:]
}

// validates the path string
func ValidatePath(path string) error {
	m, err := regexp.Match("^(/[a-z0-9]+(/[a-z0-9]+)*)?$", []byte(path))
	if err != nil || !m {
		return fmt.Errorf("path string is not in the format [/path1/path2]")
	}
	return nil
}
