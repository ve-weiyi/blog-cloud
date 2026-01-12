//go:build windows
// +build windows

package color

import (
	"fmt"
	"math/rand"
	"strconv"
)

var _ = RandomColor()

// RandomColor generates a random color.
func RandomColor() string {
	return fmt.Sprintf("#%s", strconv.FormatInt(int64(rand.Intn(16777216)), 16))
}

// Yellow ...
func Yellow(msg string) string {
	return fmt.Sprintf("%s", msg)
}

// Red ...
func Red(msg string) string {
	return fmt.Sprintf("%s", msg)
}

// Blue ...
func Blue(msg string) string {
	return fmt.Sprintf("%s", msg)
}

// Green ...
func Green(msg string) string {
	return fmt.Sprintf("%s", msg)
}
