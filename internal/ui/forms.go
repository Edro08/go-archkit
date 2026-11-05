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
)

// FormExtChange presenta el formulario interactivo para cambio de extensiones.
func FormExtChange(defaults core.ExtChangeOptions) (core.ExtChangeOptions, bool, error) {
	fmt.Print(RenderHeader("Cambiar Extensiones"))

	dir := defaults.Directory
	if dir == "" {
		dir = "."
	}

	currentExt := defaults.CurrentExt
	newExt := defaults.NewExt
	recursive := defaults.Recursive

	// 1. Preguntar por directorio
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Carpeta").
				Description("Ruta relativa o absoluta (ej. . para actual)").
				Value(&dir).
				Validate(func(s string) error {
					s = strings.TrimSpace(s)
					if s == "" {
						return errors.New("la carpeta no puede estar vacía")
					}
					info, err := os.Stat(s)
					if err != nil || !info.IsDir() {
						return fmt.Errorf("la carpeta %q no existe", s)
					}
					return nil
				}),
		),
	).WithTheme(CustomTheme()).Run()

	if err != nil {
		return defaults, true, err
	}

	absDir, _ := filepath.Abs(dir)

	// 2. Escanear extensiones para sugerir al usuario
	stats, _ := core.ScanExtensions(absDir, recursive)
	var extOptions []huh.Option[string]
	for _, s := range stats {
		if s.Extension != "(sin extensión)" {
			extOptions = append(extOptions, huh.NewOption(
				fmt.Sprintf("%s (%d archivos)", s.Extension, s.Count),
				s.Extension,
			))
		}
	}
	extOptions = append(extOptions, huh.NewOption("Otra (escribir manualmente)", "__manual__"))

	selectedExt := "__manual__"
	if len(stats) > 0 && stats[0].Extension != "(sin extensión)" {
		selectedExt = stats[0].Extension
	}

	if len(extOptions) > 1 {
		_ = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Extensión actual detectada").
					Options(extOptions...).
					Value(&selectedExt),
			),
		).WithTheme(CustomTheme()).Run()
	}

	if selectedExt == "__manual__" {
		_ = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Extensión actual").
					Placeholder(".jpeg").
					Value(&currentExt),
			),
		).WithTheme(CustomTheme()).Run()
	} else {
		currentExt = selectedExt
	}

	// 3. Nueva extensión y modo recursivo
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Nueva extensión").
				Placeholder(".jpg").
				Value(&newExt).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("debes indicar la nueva extensión")
					}
					return nil
				}),
			huh.NewConfirm().
				Title("¿Incluir subcarpetas?").
				Value(&recursive),
		),
	).WithTheme(CustomTheme())

	if err := form.Run(); err != nil {
		return defaults, true, err
	}

	return core.ExtChangeOptions{
		Directory:  absDir,
		CurrentExt: currentExt,
		NewExt:     newExt,
		Recursive:  recursive,
	}, false, nil
}

// FormNumbering presenta el formulario interactivo para numerar archivos en serie.
func FormNumbering(defaults core.NumbererOptions) (core.NumbererOptions, bool, error) {
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
	keepOriginal := defaults.KeepOriginalName
	inPlace := defaults.InPlace
	targetSubdir := defaults.TargetSubdir
	if targetSubdir == "" {
		targetSubdir = "numerados"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Carpeta").
				Value(&dir).
				Validate(validateDir),
			huh.NewInput().
				Title("Filtrar por extensión (opcional)").
				Placeholder("ej. .jpg o dejar vacío para todos").
				Value(&ext),
			huh.NewInput().
				Title("Número inicial").
				Placeholder("1").
				Value(&startStr).
				Validate(func(s string) error {
					_, ok := new(big.Int).SetString(strings.TrimSpace(s), 10)
					if !ok {
						return errors.New("debe ser un número entero válido")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Prefijo (opcional)").
				Placeholder("ej. doc_").
				Value(&prefix),
			huh.NewInput().
				Title("Sufijo (opcional)").
				Placeholder("ej. _final").
				Value(&suffix),
			huh.NewConfirm().
				Title("¿Conservar nombre original tras el número? (ej. 001_foto.jpg)").
				Value(&keepOriginal),
			huh.NewConfirm().
				Title("¿Renombrar en la misma carpeta (in-place)?").
				Description("Si no, se guardarán en una subcarpeta").
				Value(&inPlace),
		),
	).WithTheme(CustomTheme())

	if err := form.Run(); err != nil {
		return defaults, true, err
	}

	if !inPlace {
		_ = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Nombre de subcarpeta destino").
					Value(&targetSubdir),
			),
		).WithTheme(CustomTheme()).Run()
	}

	absDir, _ := filepath.Abs(dir)
	startNum, _ := new(big.Int).SetString(strings.TrimSpace(startStr), 10)

	return core.NumbererOptions{
		Directory:        absDir,
		Extension:        ext,
		Start:            startNum,
		Width:            defaults.Width,
		Prefix:           prefix,
		Suffix:           suffix,
		KeepOriginalName: keepOriginal,
		InPlace:          inPlace,
		TargetSubdir:     targetSubdir,
	}, false, nil
}

// FormRenaming presenta el formulario interactivo para buscar y reemplazar en nombres.
func FormRenaming(defaults core.RenamerOptions) (core.RenamerOptions, bool, error) {
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
				Title("Carpeta").
				Value(&dir).
				Validate(validateDir),
			huh.NewInput().
				Title("Filtrar por extensión (opcional)").
				Placeholder("ej. .png o vacío para todos").
				Value(&ext),
			huh.NewInput().
				Title("Texto a buscar").
				Placeholder("ej. copia").
				Value(&search),
			huh.NewInput().
				Title("Reemplazar por").
				Placeholder("ej. original (o dejar vacío para borrar)").
				Value(&replace),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Agregar prefijo (opcional)").
				Placeholder("ej. 2026_").
				Value(&prefix),
			huh.NewInput().
				Title("Agregar sufijo (opcional)").
				Placeholder("ej. _v2").
				Value(&suffix),
			huh.NewConfirm().
				Title("¿Ignorar mayúsculas y minúsculas al buscar?").
				Value(&ignoreCase),
			huh.NewConfirm().
				Title("¿Incluir subcarpetas recursivamente?").
				Value(&recursive),
		),
	).WithTheme(CustomTheme())

	if err := form.Run(); err != nil {
		return defaults, true, err
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
				Title("Carpeta").
				Value(&dir).
				Validate(validateDir),
			huh.NewInput().
				Title("Filtrar por extensión (opcional)").
				Placeholder("ej. .md o vacío para todos").
				Value(&ext),
			huh.NewSelect[string]().
				Title("Estilo de normalización").
				Options(styleOptions...).
				Value(&styleStr),
			huh.NewConfirm().
				Title("¿Incluir subcarpetas recursivamente?").
				Value(&recursive),
		),
	).WithTheme(CustomTheme())

	if err := form.Run(); err != nil {
		return defaults, true, err
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
