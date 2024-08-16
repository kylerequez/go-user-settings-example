package utils

import "strings"

func DisplayMessage(msg string) string {
	return strings.Title(strings.ToLower(msg))
}
