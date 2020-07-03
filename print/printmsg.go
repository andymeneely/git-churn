package print

// Ref.: https://github.com/src-d/go-git/blob/master/_examples/common.go
import (
	"fmt"
	"os"
	"strings"
)

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		PrintInCyan("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// PrintInBlue should be used to describe the example commands that are about to run.
func PrintInBlue(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// PrintInCyan should be used to display a warning
func PrintInCyan(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func PrintInGreen(format string, args ...interface{}) {
	fmt.Printf("\x1b[32;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
func PrintInYellow(format string, args ...interface{}) {
	fmt.Printf("\x1b[33;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
func PrintInPink(format string, args ...interface{}) {
	fmt.Printf("\x1b[35;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
func PrintInGrey(format string, args ...interface{}) {
	fmt.Printf("\x1b[37;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
