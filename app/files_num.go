package app

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

type FileNumberer interface {
	Number() error
}

type FileNumber struct {
	reader *bufio.Reader
}

type numberChange struct {
	oldName string
	newName string
}

func NewFileNumber(reader *bufio.Reader) *FileNumber {
	return &FileNumber{reader: reader}
}

func (n *FileNumber) Number() error {
	directory, err := chooseDirectory(n.reader)
	if err != nil {
		return err
	}
	extension, err := readExtension(n.reader, "Extensión de entrada: ")
	if err != nil {
		return err
	}

	var start *big.Int
	for {
		answer, err := readAnswer(n.reader, "Número inicial (entero mayor o igual a cero): ")
		if err != nil {
			return err
		}
		answer = strings.TrimSpace(answer)
		if answer != "" && !strings.ContainsFunc(answer, func(r rune) bool { return r < '0' || r > '9' }) {
			start, _ = new(big.Int).SetString(answer, 10)
			break
		}
		fmt.Println("Ingrese un número entero mayor o igual a cero.")
	}

	changes, width, err := n.findFiles(directory, extension, start)
	if err != nil {
		return err
	}
	if len(changes) == 0 {
		fmt.Printf("No se encontraron archivos %s en %s.\n", extension, directory)
		return nil
	}

	output := filepath.Join(directory, "temp")
	if info, err := os.Lstat(output); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("no se puede usar %q: ya existe un archivo con ese nombre", output)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("no se puede consultar la carpeta de salida: %w", err)
	}

	fmt.Printf("Archivos encontrados: %d. Numeración de %d dígitos.\n", len(changes), width)
	fmt.Printf("Orden y nombres de salida (%s):\n", output)
	for _, item := range changes {
		fmt.Printf("  %s -> %s\n", item.oldName, item.newName)
	}

	ok, err := confirm(n.reader, "¿Confirmar copia?")
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println("Proceso cancelado.")
		return nil
	}
	if err := os.MkdirAll(output, 0755); err != nil {
		return fmt.Errorf("no se pudo crear la carpeta de salida: %w", err)
	}
	n.applyChanges(directory, output, changes)
	return nil
}

func (n *FileNumber) findFiles(directory, extension string, start *big.Int) ([]numberChange, int, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, 0, fmt.Errorf("no se pudieron consultar los archivos: %w", err)
	}
	var names []string
	for _, entry := range entries {
		if entry.Type().IsRegular() && strings.EqualFold(filepath.Ext(entry.Name()), extension) {
			names = append(names, entry.Name())
		}
	}
	if len(names) == 0 {
		return nil, 0, nil
	}
	sort.Slice(names, func(i, j int) bool { return naturalLess(names[i], names[j]) })
	last := new(big.Int).Add(start, big.NewInt(int64(len(names)-1)))
	width := max(3, len(last.String()))
	current := new(big.Int).Set(start)
	changes := make([]numberChange, 0, len(names))
	for _, name := range names {
		number := current.String()
		newName := strings.Repeat("0", width-len(number)) + number + filepath.Ext(name)
		changes = append(changes, numberChange{oldName: name, newName: newName})
		current.Add(current, big.NewInt(1))
	}
	return changes, width, nil
}

func (n *FileNumber) applyChanges(directory, output string, changes []numberChange) {
	copied, skipped, failed := 0, 0, 0
	for _, item := range changes {
		source := filepath.Join(directory, item.oldName)
		destination := filepath.Join(output, item.newName)
		if err := copyFileExclusive(source, destination); err != nil {
			if os.IsExist(err) {
				fmt.Printf("Se omite %q: ya existe %q.\n", item.oldName, item.newName)
				skipped++
			} else {
				fmt.Printf("No se pudo copiar %q: %v\n", item.oldName, err)
				failed++
			}
			continue
		}
		fmt.Printf("%s -> %s\n", item.oldName, item.newName)
		copied++
	}
	fmt.Printf("Finalizado. Copiados: %d; omitidos: %d; errores: %d.\n", copied, skipped, failed)
}

func copyFileExclusive(source, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	output, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	_, err = io.Copy(output, input)
	if closeErr := output.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(destination)
	}
	return err
}

func naturalLess(a, b string) bool {
	x, y := []rune(a), []rune(b)
	i, j, zeroTie := 0, 0, 0
	digit := func(r rune) bool { return r >= '0' && r <= '9' }
	for i < len(x) && j < len(y) {
		if digit(x[i]) && digit(y[j]) {
			startI, startJ := i, j
			for i < len(x) && digit(x[i]) {
				i++
			}
			for j < len(y) && digit(y[j]) {
				j++
			}
			valueI, valueJ := startI, startJ
			for valueI < i && x[valueI] == '0' {
				valueI++
			}
			for valueJ < j && y[valueJ] == '0' {
				valueJ++
			}
			if i-valueI != j-valueJ {
				return i-valueI < j-valueJ
			}
			if cmp := strings.Compare(string(x[valueI:i]), string(y[valueJ:j])); cmp != 0 {
				return cmp < 0
			}
			if zeroTie == 0 {
				zeroTie = (i - startI) - (j - startJ)
			}
		} else {
			left, right := unicode.ToLower(x[i]), unicode.ToLower(y[j])
			if left != right {
				return left < right
			}
			i++
			j++
		}
	}
	if i != len(x) || j != len(y) {
		return len(x)-i < len(y)-j
	}
	if zeroTie != 0 {
		return zeroTie < 0
	}
	return a < b
}
