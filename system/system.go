package system

import "fmt"

func CheckString(v interface{}) (string, error) {
	switch v.(type) {
	case string:
		return v.(string), nil
	default:
		return "", fmt.Errorf("some errors", v)
	}
}
