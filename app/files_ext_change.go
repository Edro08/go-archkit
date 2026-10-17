package app

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileExtChanger interface {
	Change() error
}

type FileExtChange struct {
	reader *bufio.Reader
}

type fileChange struct {
	oldName string
	newName string
}

func NewChangeExt(reader *bufio.Reader) *FileExtChange {
	return &FileExtChange{reader: reader}
}

func (c *FileExtChange) Change() error {
	directory, err := chooseDirectory(c.reader)
	if err != nil {
		return err
	}

	currentExt, err := readExtension(c.reader, "Extensión actual: ")
	if err != nil {
		return err
	}

	newExt, err := readExtension(c.reader, "Extensión nueva: ")
	if err != nil {
		return err
	}

	if strings.EqualFold(currentExt, newExt) {
		fmt.Println("Las extensiones son iguales. No hay cambios que realizar.")
		return nil
	}

	changes, err := c.findFiles(directory, currentExt, newExt)
	if err != nil {
		return err
	}

	if len(changes) == 0 {
		fmt.Println("No se encontraron archivos para modificar.")
		return nil
	}

	fmt.Printf("Se encontraron %d archivos:\n", len(changes))
	for _, item := range changes {
		fmt.Printf("  %s -> %s\n", item.oldName, item.newName)
	}

	fmt.Println("Solo cambiará la extensión; el contenido no se convertirá.")

	ok, err := confirm(c.reader, fmt.Sprintf("¿Cambiar a %s?", newExt))
	if err != nil {
		return err
	}

	if !ok {
		fmt.Println("Proceso cancelado.")
		return nil
	}

	c.applyChanges(directory, changes)
	return nil
}

func (c *FileExtChange) findFiles(directory, currentExt, newExt string) ([]fileChange, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron consultar los archivos: %w", err)
	}

	var changes []fileChange

	for _, entry := range entries {
		if !entry.Type().IsRegular() ||
			!strings.EqualFold(filepath.Ext(entry.Name()), currentExt) {
			continue
		}

		oldName := entry.Name()
		newName := strings.TrimSuffix(oldName, filepath.Ext(oldName)) + newExt

		changes = append(changes, fileChange{
			oldName: oldName,
			newName: newName,
		})
	}

	return changes, nil
}

func (c *FileExtChange) applyChanges(directory string, changes []fileChange) {
	renamed, skipped, failed := 0, 0, 0

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
