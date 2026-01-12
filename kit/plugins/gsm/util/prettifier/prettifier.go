package prettifier

import (
	"strings"

	snake "github.com/azer/snakecase"
	"github.com/rainycape/unidecode"
)

// trim, snake case, no space, no dot
func PrettifyKey(s string) string {
	return replaceDot(replaceSpace(snake.SnakeCase(strings.TrimSpace(s))))
}

func replaceDot(s string) string {
	return strings.Replace(s, ".", "_", -1)
}

func replaceSpace(s string) string {
	return strings.Replace(s, " ", "_", -1)
}

// trim, format list
func PrettifyValue(s string) interface{} {
	// replace non-ascii chars
	s = strings.TrimSpace(unidecode.Unidecode(s))

	// some case, there's 4 space, let's treat it as a list
	s = strings.Replace(s, "    ", "\n-", -1)

	// convert unordered list into slice
	listString := strings.Split(s, "\n")
	if len(listString) == 1 {
		return listString[0]
	}

	for key, value := range listString {
		listString[key] = strings.Replace(value, "- ", "", -1)
	}

	return listString
}
