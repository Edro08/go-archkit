package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ExtChangeOptions define los parámetros para el cambio masivo de extensiones.
type ExtChangeOptions struct {
	Directory  string
	CurrentExt string // Ej: ".jpeg" o "jpeg"
	NewExt     string // Ej: ".jpg" o "jpg"
	Recursive  bool
}

// NormalizeExtension asegura que la extensión empiece por punto y esté en minúsculas.
func NormalizeExtension(ext string) string {
	ext = strings.TrimSpace(ext)
	if ext == "" {
		return ""
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return ext
}

// ExtensionStat almacena una extensión encontrada y su recuento de archivos.
type ExtensionStat struct {
	Extension string
	Count     int
}

// ScanExtensions escanea un directorio y retorna las extensiones presentes con su frecuencia.
func ScanExtensions(dir string, recursive bool) ([]ExtensionStat, error) {
	counts := make(map[string]int)

	walkFn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if !recursive && path != dir {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Type().IsRegular() {
			ext := strings.ToLower(filepath.Ext(d.Name()))
			if ext != "" {
				counts[ext]++
			} else {
				counts["(sin extensión)"]++
			}
		}
		return nil
	}

	if err := filepath.WalkDir(dir, walkFn); err != nil {
		return nil, fmt.Errorf("error al escanear extensiones: %w", err)
	}

	stats := make([]ExtensionStat, 0, len(counts))
	for ext, count := range counts {
		stats = append(stats, ExtensionStat{Extension: ext, Count: count})
	}

	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Count == stats[j].Count {
			return stats[i].Extension < stats[j].Extension
		}
		return stats[i].Count > stats[j].Count
	})

	return stats, nil
}

// PlanExtChange calcula el plan de modificación de extensiones sin tocar el disco.
func PlanExtChange(opts ExtChangeOptions) (*FilePlan, error) {
	dir, err := filepath.Abs(opts.Directory)
	if err != nil {
		return nil, fmt.Errorf("ruta de directorio inválida: %w", err)
	}

	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("el directorio %q no existe o no es válido", dir)
	}

	currentExt := NormalizeExtension(opts.CurrentExt)
	newExt := NormalizeExtension(opts.NewExt)

	if strings.EqualFold(currentExt, newExt) {
		return nil, fmt.Errorf("la extensión actual y la nueva son idénticas (%s)", newExt)
	}

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

		if strings.EqualFold(ext, currentExt) {
			fileInfo, err := d.Info()
			if err != nil {
				return err
			}

			baseNameWithoutExt := strings.TrimSuffix(name, ext)
			newName := baseNameWithoutExt + newExt
			targetPath := filepath.Join(filepath.Dir(path), newName)

			plan.Add(FileChange{
				SourcePath: path,
				TargetPath: targetPath,
				OldName:    name,
				NewName:    newName,
				Size:       fileInfo.Size(),
				IsDir:      false,
			})
		}
		return nil
	}

	if err := filepath.WalkDir(dir, walkFn); err != nil {
		return nil, fmt.Errorf("error al buscar archivos: %w", err)
	}

	return plan, nil
}
