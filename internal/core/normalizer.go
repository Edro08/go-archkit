package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// CaseStyle define el estilo de capitalización y delimitadores de palabras.
type CaseStyle string

const (
	CaseKebab  CaseStyle = "kebab-case"  // mi-nombre-de-archivo
	CaseSnake  CaseStyle = "snake_case"  // mi_nombre_de_archivo
	CaseCamel  CaseStyle = "camelCase"   // miNombreDeArchivo
	CasePascal CaseStyle = "PascalCase"  // MiNombreDeArchivo
	CaseLower  CaseStyle = "lowercase"   // minombresuelto
	CaseUpper  CaseStyle = "UPPERCASE"   // MINOMBRESUELTO
)

// NormalizerOptions define los parámetros para normalizar nombres de archivo.
type NormalizerOptions struct {
	Directory string
	Extension string
	Style     CaseStyle
	Recursive bool
}

// PlanNormalize calcula el plan para convertir los nombres al estilo elegido.
func PlanNormalize(opts NormalizerOptions) (*FilePlan, error) {
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
		normalizedBase := ConvertCase(baseName, opts.Style)
		if normalizedBase == "" {
			normalizedBase = baseName
		}

		newName := normalizedBase + ext
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
		return nil, fmt.Errorf("error al escanear archivos para normalizar: %w", err)
	}

	return plan, nil
}

// ConvertCase transforma un string al estilo especificado.
func ConvertCase(input string, style CaseStyle) string {
	words := splitWords(input)
	if len(words) == 0 {
		return input
	}

	switch style {
	case CaseKebab:
		lowers := make([]string, len(words))
		for i, w := range words {
			lowers[i] = strings.ToLower(w)
		}
		return strings.Join(lowers, "-")

	case CaseSnake:
		lowers := make([]string, len(words))
		for i, w := range words {
			lowers[i] = strings.ToLower(w)
		}
		return strings.Join(lowers, "_")

	case CaseCamel:
		var b strings.Builder
		for i, w := range words {
			if i == 0 {
				b.WriteString(strings.ToLower(w))
			} else {
				b.WriteString(capitalize(w))
			}
		}
		return b.String()

	case CasePascal:
		var b strings.Builder
		for _, w := range words {
			b.WriteString(capitalize(w))
		}
		return b.String()

	case CaseLower:
		return strings.ToLower(input)

	case CaseUpper:
		return strings.ToUpper(input)

	default:
		return input
	}
}

func splitWords(input string) []string {
	var words []string
	var current strings.Builder

	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == '-' || r == '_' || r == '.' || unicode.IsSpace(r) {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
			continue
		}

		// Detección de camelCase: si es mayúscula y la anterior era minúscula
		if unicode.IsUpper(r) && i > 0 && unicode.IsLower(runes[i-1]) {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
		}

		current.WriteRune(r)
	}

	if current.Len() > 0 {
		words = append(words, current.String())
	}

	return words
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}
