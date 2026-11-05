package core

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// NumbererOptions define los parámetros para numerar archivos en serie.
type NumbererOptions struct {
	Directory        string
	Extension        string   // Opcional: filtrar por extensión (ej. ".png")
	Start            *big.Int // Número inicial (ej. 1)
	Width            int      // Dígitos de padding (0 para auto-calcular)
	Prefix           string   // Prefijo opcional (ej. "img_")
	Suffix           string   // Sufijo opcional antes de la extensión (ej. "_final")
	KeepOriginalName bool     // Si true: "001_original.jpg", si false: "001.jpg"
	InPlace          bool     // Si true: renombra en la misma carpeta; si false: a subcarpeta
	TargetSubdir     string   // Nombre de la subcarpeta si InPlace es false (default "numerados")
}

// PlanNumbering calcula el plan para numerar archivos con orden natural.
func PlanNumbering(opts NumbererOptions) (*FilePlan, error) {
	dir, err := filepath.Abs(opts.Directory)
	if err != nil {
		return nil, fmt.Errorf("ruta de directorio inválida: %w", err)
	}

	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("el directorio %q no existe o no es válido", dir)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el directorio: %w", err)
	}

	extFilter := NormalizeExtension(opts.Extension)

	// Filtrar y recolectar archivos coincidentes
	type targetEntry struct {
		name string
		size int64
	}
	var matched []targetEntry

	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		if extFilter != "" && !strings.EqualFold(filepath.Ext(entry.Name()), extFilter) {
			continue
		}
		entryInfo, err := entry.Info()
		if err != nil {
			continue
		}
		matched = append(matched, targetEntry{name: entry.Name(), size: entryInfo.Size()})
	}

	if len(matched) == 0 {
		return NewFilePlan(dir), nil
	}

	// Orden natural: "2.txt" antes de "10.txt"
	sort.Slice(matched, func(i, j int) bool {
		return NaturalLess(matched[i].name, matched[j].name)
	})

	start := opts.Start
	if start == nil {
		start = big.NewInt(1)
	}

	// Calcular ancho de dígitos
	width := opts.Width
	if width <= 0 {
		totalCount := big.NewInt(int64(len(matched) - 1))
		lastNum := new(big.Int).Add(start, totalCount)
		width = max(3, len(lastNum.String()))
	}

	targetDir := dir
	if !opts.InPlace {
		sub := opts.TargetSubdir
		if strings.TrimSpace(sub) == "" {
			sub = "numerados"
		}
		targetDir = filepath.Join(dir, sub)
	}

	plan := NewFilePlan(dir)
	current := new(big.Int).Set(start)

	for _, item := range matched {
		numStr := current.String()
		padded := numStr
		if len(numStr) < width {
			padded = strings.Repeat("0", width-len(numStr)) + numStr
		}

		ext := filepath.Ext(item.name)
		rawBase := strings.TrimSuffix(item.name, ext)

		var newBase string
		if opts.KeepOriginalName {
			newBase = fmt.Sprintf("%s%s_%s%s", opts.Prefix, padded, rawBase, opts.Suffix)
		} else {
			newBase = fmt.Sprintf("%s%s%s", opts.Prefix, padded, opts.Suffix)
		}

		newName := newBase + ext
		targetPath := filepath.Join(targetDir, newName)

		plan.Add(FileChange{
			SourcePath: filepath.Join(dir, item.name),
			TargetPath: targetPath,
			OldName:    item.name,
			NewName:    newName,
			Size:       item.size,
			IsDir:      false,
		})

		current.Add(current, big.NewInt(1))
	}

	return plan, nil
}
