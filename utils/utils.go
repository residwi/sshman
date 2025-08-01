package utils

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

const (
	WarningSymbol = "\u26A0"
	SuccessSymbol = "\u2714"
	ErrorSymbol   = "\u2716"
)

func IsFileNotExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return os.IsNotExist(err)
}

func IsDirectoryNotExist(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil {
		return os.IsNotExist(err)
	}
	return !info.IsDir()
}

func PrintSuccess(message string) {
	color.Green("%s %s\n", SuccessSymbol, message)
}

func PrintError(message string) {
	color.Red("%s %s\n", ErrorSymbol, message)
}

func PrintWarning(message string) {
	color.Yellow("%s %s\n", WarningSymbol, message)
}

func PrintTable(w io.Writer, headers []string, rows [][]string) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, strings.Join(headers, "\t"))

	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}

	tw.Flush()
}

func ReplaceHomeDirWithTilde(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}
