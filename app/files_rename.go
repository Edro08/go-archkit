package app

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type FileRenamer interface {
	Rename() error
}

type FileRename struct {
	reader *bufio.Reader
}

type renameChange struct {
	oldName string
	newName string
}

type renameOptions struct {
	extension   string
	search      string
	replacement string
	prefix      string
	suffix      string
}

func NewFileRename(reader *bufio.Reader) *FileRename {
	return &FileRename{reader: reader}
}

func (r *FileRename) Rename() error {
	directory, err := chooseDirectory(r.reader)
	if err != nil {
		return err
	}
	options := renameOptions{}
	options.extension, err = readExtension(r.reader, "Extensión de los archivos: ")
	if err != nil {
		return err
	}
	options.search, err = readAnswer(r.reader, "Texto a buscar (Enter para no reemplazar): ")
	if err != nil {
		return err
	}
	if options.search != "" {
		options.replacement, err = readAnswer(r.reader, "Reemplazar por (Enter para eliminar el texto): ")
		if err != nil {
			return err
		}
	}
	options.prefix, err = readAnswer(r.reader, "Prefijo (Enter para ninguno): ")
	if err != nil {
		return err
	}
	options.suffix, err = readAnswer(r.reader, "Sufijo antes de la extensión (Enter para ninguno): ")
	if err != nil {
		return err
	}
	if options.search == "" && options.prefix == "" && options.suffix == "" {
		fmt.Println("No se indicó ningún cambio. Proceso cancelado.")
		return nil
	}

	changes, skipped, err := r.findFiles(directory, options)
	if err != nil {
		return err
	}
	if len(changes) == 0 {
		fmt.Printf("No hay archivos para renombrar. Omitidos: %d.\n", skipped)
		return nil
	}
	fmt.Printf("Cambios propuestos en %s:\n", directory)
	for _, item := range changes {
		fmt.Printf("  %s -> %s\n", item.oldName, item.newName)
	}
	fmt.Printf("Archivos por renombrar: %d; omitidos: %d.\n", len(changes), skipped)
	ok, err := confirm(r.reader, "¿Confirmar cambios?")
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println("Proceso cancelado.")
		return nil
	}
	r.applyChanges(directory, changes, skipped)
	return nil
}

func (r *FileRename) findFiles(directory string, options renameOptions) ([]renameChange, int, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, 0, fmt.Errorf("no se pudieron consultar los archivos: %w", err)
	}
	var changes []renameChange
	skipped := 0
	planned := make(map[string]bool)
	for _, entry := range entries {
		name := entry.Name()
		if !entry.Type().IsRegular() || !strings.EqualFold(filepath.Ext(name), options.extension) {
			continue
		}
		extension := filepath.Ext(name)
		base := strings.TrimSuffix(name, extension)
		if options.search != "" {
			base = strings.ReplaceAll(base, options.search, options.replacement)
		}
		newName := options.prefix + base + options.suffix + extension
		if newName == name {
			skipped++
			continue
		}
		if !validFileName(newName) {
			fmt.Printf("Se omite %q: el nombre de destino no es válido.\n", name)
			skipped++
			continue
		}
		if planned[strings.ToLower(newName)] {
			fmt.Printf("Se omite %q: el destino %q ya está previsto.\n", name, newName)
			skipped++
			continue
		}
		if _, err := os.Lstat(filepath.Join(directory, newName)); err == nil {
			fmt.Printf("Se omite %q: ya existe %q.\n", name, newName)
			skipped++
			continue
		} else if !os.IsNotExist(err) {
			return nil, skipped, fmt.Errorf("no se pudo comprobar %q: %w", newName, err)
		}
		planned[strings.ToLower(newName)] = true
		changes = append(changes, renameChange{oldName: name, newName: newName})
	}
	return changes, skipped, nil
}

func (r *FileRename) applyChanges(directory string, changes []renameChange, skipped int) {
	renamed, failed := 0, 0
	for _, item := range changes {
		source := filepath.Join(directory, item.oldName)
		destination := filepath.Join(directory, item.newName)
		if _, err := os.Lstat(destination); err == nil {
			fmt.Printf("Se omite %q: ya existe %q.\n", item.oldName, item.newName)
			skipped++
			continue
		} else if !os.IsNotExist(err) {
			fmt.Printf("No se pudo comprobar %q: %v\n", item.newName, err)
			failed++
			continue
		}
		if err := os.Rename(source, destination); err != nil {
			fmt.Printf("No se pudo renombrar %q: %v\n", item.oldName, err)
			failed++
			continue
		}
		fmt.Printf("%s -> %s\n", item.oldName, item.newName)
		renamed++
	}
	fmt.Printf("Finalizado. Renombrados: %d; omitidos: %d; errores: %d.\n", renamed, skipped, failed)
}

func validFileName(name string) bool {
	if name == "" || strings.TrimSpace(strings.TrimSuffix(name, filepath.Ext(name))) == "" ||
		strings.HasSuffix(name, ".") || strings.HasSuffix(name, " ") ||
		strings.ContainsAny(name, "\\/:*?\"<>|") || strings.ContainsFunc(name, unicode.IsControl) {
		return false
	}
	base := strings.ToUpper(strings.SplitN(name, ".", 2)[0])
	switch base {
	case "CON", "PRN", "AUX", "NUL", "CONIN$", "CONOUT$":
		return false
	}
	if len(base) == 4 && (strings.HasPrefix(base, "COM") || strings.HasPrefix(base, "LPT")) && strings.ContainsAny(base[3:], "123456789") {
		return false
	}
	return true
}
