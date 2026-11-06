package ui

import (
	"errors"
	"fmt"
	"go-archkit/internal/core"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// FormExtChange presenta el formulario interactivo para cambio de extensiones.
func FormExtChange(defaults core.ExtChangeOptions) (core.ExtChangeOptions, bool, error) {
	ClearScreen()
	fmt.Print(RenderHeader("Cambiar Extensiones"))

	dir := defaults.Directory
	if dir == "" {
		dir = "."
	}

	currentExt := defaults.CurrentExt
	newExt := defaults.NewExt
	recursive := defaults.Recursive

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Carpeta:").
				Placeholder(". (actual)").
				Value(&dir).
				Inline(true).
				Validate(validateDir),
			huh.NewInput().
				Title("Extensión actual:").
				Placeholder("ej. .jpeg").
				Value(&currentExt).
				Inline(true).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("debes indicar la extensión actual")
					}
					return nil
				}),
			huh.NewInput().
				Title("Nueva extensión:").
				Placeholder("ej. .jpg").
				Value(&newExt).
				Inline(true).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("debes indicar la nueva extensión")
					}
					return nil
				}),
			huh.NewConfirm().
				Title("¿Incluir subcarpetas?").
				Value(&recursive).
				Inline(true).
				WithButtonAlignment(lipgloss.Left),
			FooterNote(),
		),
	).WithTheme(CustomTheme()).WithKeyMap(CustomKeyMap()).WithShowHelp(false)

	if err := form.Run(); err != nil {
		return defaults, true, nil
	}

	absDir, _ := filepath.Abs(dir)

	return core.ExtChangeOptions{
		Directory:  absDir,
		CurrentExt: currentExt,
		NewExt:     newExt,
		Recursive:  recursive,
	}, false, nil
}

// FormNumbering presenta el formulario interactivo para numerar archivos en serie.
func FormNumbering(defaults core.NumbererOptions) (core.NumbererOptions, bool, error) {
	ClearScreen()
	fmt.Print(RenderHeader("Numerar Archivos en Serie"))

	dir := defaults.Directory
	if dir == "" {
		dir = "."
	}

	ext := defaults.Extension
	startStr := "1"
	if defaults.Start != nil {
		startStr = defaults.Start.String()
	}
	prefix := defaults.Prefix
	suffix := defaults.Suffix

	// Por defecto: NO conservar nombre original (false)
	keepOriginal := defaults.KeepOriginalName

	// Por defecto: NO renombrar in-place (false)
	inPlace := defaults.InPlace

	// Por defecto: subcarpeta "temp"
	targetSubdir := defaults.TargetSubdir
	if strings.TrimSpace(targetSubdir) == "" {
		targetSubdir = "temp"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Carpeta:").
				Value(&dir).
				Inline(true).
				Validate(validateDir),
			huh.NewInput().
				Title("Filtrar por extensión (opcional):").
				Placeholder("ej. .jpg o vacío para todos").
				Value(&ext).
				Inline(true),
			huh.NewInput().
				Title("Número inicial:").
				Placeholder("1").
				Value(&startStr).
				Inline(true).
				Validate(func(s string) error {
					_, ok := new(big.Int).SetString(strings.TrimSpace(s), 10)
					if !ok {
						return errors.New("debe ser un número entero válido")
					}
					return nil
				}),
			huh.NewInput().
				Title("Prefijo (opcional):").
				Placeholder("ej. doc_").
				Value(&prefix).
				Inline(true),
			huh.NewInput().
				Title("Sufijo (opcional):").
				Placeholder("ej. _final").
				Value(&suffix).
				Inline(true),
			huh.NewConfirm().
				Title("¿Conservar nombre original tras número?").
				Value(&keepOriginal).
				Inline(true).
				WithButtonAlignment(lipgloss.Left),
			huh.NewConfirm().
				Title("¿Renombrar in-place (misma carpeta)?").
				Value(&inPlace).
				Inline(true).
				WithButtonAlignment(lipgloss.Left),
			huh.NewInput().
				Title("Subcarpeta de destino:").
				Placeholder("temp (por defecto)").
				Value(&targetSubdir).
				Inline(true),
			FooterNote(),
		),
	).WithTheme(CustomTheme()).WithKeyMap(CustomKeyMap()).WithShowHelp(false)

	if err := form.Run(); err != nil {
		return defaults, true, nil
	}

	absDir, _ := filepath.Abs(dir)
	startClean := strings.TrimSpace(startStr)
	startNum, _ := new(big.Int).SetString(startClean, 10)
	width := len(startClean)

	cleanSubdir := strings.TrimSpace(targetSubdir)
	if cleanSubdir == "" {
		cleanSubdir = "temp"
	}

	return core.NumbererOptions{
		Directory:        absDir,
		Extension:        ext,
		Start:            startNum,
		Width:            width,
		Prefix:           prefix,
		Suffix:           suffix,
		KeepOriginalName: keepOriginal,
		InPlace:          inPlace,
		TargetSubdir:     cleanSubdir,
	}, false, nil
}

// FormRenaming presenta el formulario interactivo para buscar y reemplazar en nombres.
func FormRenaming(defaults core.RenamerOptions) (core.RenamerOptions, bool, error) {
	ClearScreen()
	fmt.Print(RenderHeader("Buscar y Reemplazar"))

	dir := defaults.Directory
	if dir == "" {
		dir = "."
	}

	ext := defaults.Extension
	search := defaults.Search
	replace := defaults.Replace
	prefix := defaults.Prefix
	suffix := defaults.Suffix
	ignoreCase := defaults.IgnoreCase
	recursive := defaults.Recursive

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Carpeta:").
				Value(&dir).
				Inline(true).
				Validate(validateDir),
			huh.NewInput().
				Title("Filtrar por extensión (opcional):").
				Placeholder("ej. .png o vacío para todos").
				Value(&ext).
				Inline(true),
			huh.NewInput().
				Title("Texto a buscar:").
				Placeholder("ej. copia").
				Value(&search).
				Inline(true),
			huh.NewInput().
				Title("Reemplazar por:").
				Placeholder("ej. original (o vacío para borrar)").
				Value(&replace).
				Inline(true),
			huh.NewInput().
				Title("Prefijo (opcional):").
				Placeholder("ej. 2026_").
				Value(&prefix).
				Inline(true),
			huh.NewInput().
				Title("Sufijo (opcional):").
				Placeholder("ej. _v2").
				Value(&suffix).
				Inline(true),
			huh.NewConfirm().
				Title("¿Ignorar mayúsculas y minúsculas?").
				Value(&ignoreCase).
				Inline(true).
				WithButtonAlignment(lipgloss.Left),
			huh.NewConfirm().
				Title("¿Incluir subcarpetas?").
				Value(&recursive).
				Inline(true).
				WithButtonAlignment(lipgloss.Left),
			FooterNote(),
		),
	).WithTheme(CustomTheme()).WithKeyMap(CustomKeyMap()).WithShowHelp(false)

	if err := form.Run(); err != nil {
		return defaults, true, nil
	}

	absDir, _ := filepath.Abs(dir)

	return core.RenamerOptions{
		Directory:  absDir,
		Extension:  ext,
		Search:     search,
		Replace:    replace,
		Prefix:     prefix,
		Suffix:     suffix,
		IgnoreCase: ignoreCase,
		Recursive:  recursive,
	}, false, nil
}

// FormNormalizing presenta el formulario interactivo para normalizar a kebab-case, snake_case, etc.
func FormNormalizing(defaults core.NormalizerOptions) (core.NormalizerOptions, bool, error) {
	ClearScreen()
	fmt.Print(RenderHeader("Normalizar Nombres de Archivo"))

	dir := defaults.Directory
	if dir == "" {
		dir = "."
	}

	ext := defaults.Extension
	styleStr := string(defaults.Style)
	if styleStr == "" {
		styleStr = string(core.CaseKebab)
	}
	recursive := defaults.Recursive

	styleOptions := []huh.Option[string]{
		huh.NewOption("kebab-case (ej. mi-archivo-de-fotos.jpg)", string(core.CaseKebab)),
		huh.NewOption("snake_case (ej. mi_archivo_de_fotos.jpg)", string(core.CaseSnake)),
		huh.NewOption("camelCase (ej. miArchivoDeFotos.jpg)", string(core.CaseCamel)),
		huh.NewOption("PascalCase (ej. MiArchivoDeFotos.jpg)", string(core.CasePascal)),
		huh.NewOption("lowercase (ej. miarchivodefotos.jpg)", string(core.CaseLower)),
		huh.NewOption("UPPERCASE (ej. MIARCHIVODEFOTOS.JPG)", string(core.CaseUpper)),
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Carpeta:").
				Value(&dir).
				Inline(true).
				Validate(validateDir),
			huh.NewInput().
				Title("Filtrar por extensión (opcional):").
				Placeholder("ej. .md o vacío para todos").
				Value(&ext).
				Inline(true),
			huh.NewSelect[string]().
				Title("Estilo de normalización:").
				Options(styleOptions...).
				Value(&styleStr),
			huh.NewConfirm().
				Title("¿Incluir subcarpetas?").
				Value(&recursive).
				Inline(true).
				WithButtonAlignment(lipgloss.Left),
			FooterNote(),
		),
	).WithTheme(CustomTheme()).WithKeyMap(CustomKeyMap()).WithShowHelp(false)

	if err := form.Run(); err != nil {
		return defaults, true, nil
	}

	absDir, _ := filepath.Abs(dir)

	return core.NormalizerOptions{
		Directory: absDir,
		Extension: ext,
		Style:     core.CaseStyle(styleStr),
		Recursive: recursive,
	}, false, nil
}

func validateDir(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return errors.New("la carpeta no puede estar vacía")
	}
	info, err := os.Stat(s)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("la carpeta %q no existe", s)
	}
	return nil
}

func validateInt(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	_, err := strconv.Atoi(s)
	if err != nil {
		return errors.New("debe ser un número entero")
	}
	return nil
}
