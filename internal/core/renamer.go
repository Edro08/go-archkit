package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// RenamerOptions define los parámetros para renombrar archivos.
type RenamerOptions struct {
	Directory  string
	Extension  string // Opcional: filtrar por extensión
	Search     string // Texto a buscar
	Replace    string // Texto con el que reemplazar
	Prefix     string // Prefijo que agregar al inicio del nombre
	Suffix     string // Sufijo que agregar al final del nombre (antes de la extensión)
	IgnoreCase bool   // Ignorar mayúsculas y minúsculas en la búsqueda
	Recursive  bool   // Procesar subdirectorios
}

// PlanRename calcula el plan de renombrado aplicando búsquedas, reemplazos y prefijos/sufijos.
func PlanRename(opts RenamerOptions) (*FilePlan, error) {
	dir, err := filepath.Abs(opts.Directory)
	if err != nil {
		return nil, fmt.Errorf("ruta de directorio inválida: %w", err)
	}

	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("el directorio %q no existe o no es válido", dir)
	}

	extFilter := NormalizeExtension(opts.Extension)
	plan := NewFilePlan(dir)

	walkFn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if !opts.Recursive && path != dir {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}

		name := d.Name()
		ext := filepath.Ext(name)

		if extFilter != "" && !strings.EqualFold(ext, extFilter) {
			return nil
		}

		baseName := strings.TrimSuffix(name, ext)
		newBaseName := baseName

		// Aplicar búsqueda y reemplazo
		if opts.Search != "" {
			if opts.IgnoreCase {
				newBaseName = replaceCaseInsensitive(newBaseName, opts.Search, opts.Replace)
			} else {
				newBaseName = strings.ReplaceAll(newBaseName, opts.Search, opts.Replace)
			}
		}

		// Aplicar prefijo y sufijo
		if opts.Prefix != "" || opts.Suffix != "" {
			newBaseName = opts.Prefix + newBaseName + opts.Suffix
		}

		newName := newBaseName + ext
		targetPath := filepath.Join(filepath.Dir(path), newName)

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		change := FileChange{
			SourcePath: path,
			TargetPath: targetPath,
			OldName:    name,
			NewName:    newName,
			Size:       fileInfo.Size(),
			IsDir:      false,
		}

		if change.IsModified() {
			plan.Add(change)
		}

		return nil
	}

	if err := filepath.WalkDir(dir, walkFn); err != nil {
		return nil, fmt.Errorf("error al escanear archivos para renombrar: %w", err)
	}

	return plan, nil
}

func replaceCaseInsensitive(input, search, replace string) string {
	if search == "" {
		return input
	}
	lowerInput := strings.ToLower(input)
	lowerSearch := strings.ToLower(search)

	var builder strings.Builder
	start := 0
	for {
		idx := strings.Index(lowerInput[start:], lowerSearch)
		if idx == -1 {
			builder.WriteString(input[start:])
			break
		}
		realIdx := start + idx
		builder.WriteString(input[start:realIdx])
		builder.WriteString(replace)
		start = realIdx + len(search)
	}
	return builder.String()
}
