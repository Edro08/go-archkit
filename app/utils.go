package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func readAnswer(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	answer, err := reader.ReadString('\n')
	if err != nil && (err != io.EOF || answer == "") {
		return "", err
	}
	return strings.TrimRight(answer, "\r\n"), nil
}

func chooseDirectory(reader *bufio.Reader) (string, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		path, err := readAnswer(reader, fmt.Sprintf("Carpeta (Enter para %s): ", current))
		if err != nil {
			return "", err
		}
		path = strings.TrimSpace(path)
		if path == "" {
			path = current
		}
		path, err = filepath.Abs(path)
		if err == nil {
			var info os.FileInfo
			info, err = os.Stat(path)
			if err == nil && info.IsDir() {
				return path, nil
			}
		}
		fmt.Println("Ingrese una carpeta existente.")
	}
}

func readExtension(reader *bufio.Reader, prompt string) (string, error) {
	for {
		extension, err := readAnswer(reader, prompt)
		if err != nil {
			return "", err
		}
		extension = strings.TrimPrefix(strings.TrimSpace(extension), ".")
		if extension != "" && !strings.ContainsAny(extension, ".\\/:*?\"<>|") &&
			!strings.ContainsFunc(extension, func(r rune) bool { return unicode.IsSpace(r) || unicode.IsControl(r) }) {
			return "." + extension, nil
		}
		fmt.Println("Ingrese una extension valida, con o sin punto.")
	}
}

func confirm(reader *bufio.Reader, prompt string) (bool, error) {
	answer, err := readAnswer(reader, prompt+" (S/N): ")
	if err != nil {
		return false, err
	}
	return strings.EqualFold(strings.TrimSpace(answer), "S"), nil
}
