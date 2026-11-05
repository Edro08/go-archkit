package ui

import (
	"fmt"
	"go-archkit/internal/core"
	"os"

	"github.com/charmbracelet/huh"
)

// Run inicia el bucle interactivo principal de Archkit.
func Run() {
	for {
		fmt.Print("\033[2J\033[H") // Limpiar pantalla en terminal
		fmt.Print(RenderHeader(""))

		var choice string
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("¿Qué deseas hacer hoy?").
					Options(
						huh.NewOption("Cambiar extensiones (.jpeg ➜ .jpg)", "ext"),
						huh.NewOption("Numerar archivos en serie (001, 002...)", "num"),
						huh.NewOption("Buscar y reemplazar / Prefijos / Sufijos", "rename"),
						huh.NewOption("Normalizar formato (kebab, snake, camel)", "norm"),
						huh.NewOption("Salir", "exit"),
					).
					Value(&choice),
			),
		).WithTheme(CustomTheme()).Run()

		if err != nil || choice == "exit" {
			fmt.Println("\n  " + styleSubtitle.Render("Hasta pronto."))
			return
		}

		switch choice {
		case "ext":
			handleExtChange()
		case "num":
			handleNumbering()
		case "rename":
			handleRenaming()
		case "norm":
			handleNormalizing()
		}
	}
}

func handleExtChange() {
	opts := core.ExtChangeOptions{Directory: "."}
	for {
		newOpts, cancelled, err := FormExtChange(opts)
		if err != nil || cancelled {
			return
		}
		opts = newOpts

		plan, err := core.PlanExtChange(opts)
		if err != nil {
			showError(err)
			continue
		}

		action, err := ShowPreview(plan)
		if err != nil || action == ActionCancel {
			return
		}
		if action == ActionEdit {
			continue
		}

		// Aplicar cambios
		res, err := core.Execute(plan, core.ModeRename)
		if err != nil {
			showError(err)
			return
		}
		showSuccess(res.Processed, len(res.Errors))
		return
	}
}

func handleNumbering() {
	opts := core.NumbererOptions{Directory: ".", InPlace: true, TargetSubdir: "numerados"}
	for {
		newOpts, cancelled, err := FormNumbering(opts)
		if err != nil || cancelled {
			return
		}
		opts = newOpts

		plan, err := core.PlanNumbering(opts)
		if err != nil {
			showError(err)
			continue
		}

		action, err := ShowPreview(plan)
		if err != nil || action == ActionCancel {
			return
		}
		if action == ActionEdit {
			continue
		}

		mode := core.ModeRename
		if !opts.InPlace {
			mode = core.ModeCopy
		}

		res, err := core.Execute(plan, mode)
		if err != nil {
			showError(err)
			return
		}
		showSuccess(res.Processed, len(res.Errors))
		return
	}
}

func handleRenaming() {
	opts := core.RenamerOptions{Directory: "."}
	for {
		newOpts, cancelled, err := FormRenaming(opts)
		if err != nil || cancelled {
			return
		}
		opts = newOpts

		plan, err := core.PlanRename(opts)
		if err != nil {
			showError(err)
			continue
		}

		action, err := ShowPreview(plan)
		if err != nil || action == ActionCancel {
			return
		}
		if action == ActionEdit {
			continue
		}

		res, err := core.Execute(plan, core.ModeRename)
		if err != nil {
			showError(err)
			return
		}
		showSuccess(res.Processed, len(res.Errors))
		return
	}
}

func handleNormalizing() {
	opts := core.NormalizerOptions{Directory: ".", Style: core.CaseKebab}
	for {
		newOpts, cancelled, err := FormNormalizing(opts)
		if err != nil || cancelled {
			return
		}
		opts = newOpts

		plan, err := core.PlanNormalize(opts)
		if err != nil {
			showError(err)
			continue
		}

		action, err := ShowPreview(plan)
		if err != nil || action == ActionCancel {
			return
		}
		if action == ActionEdit {
			continue
		}

		res, err := core.Execute(plan, core.ModeRename)
		if err != nil {
			showError(err)
			return
		}
		showSuccess(res.Processed, len(res.Errors))
		return
	}
}

func showSuccess(processed, errorCount int) {
	fmt.Println()
	if errorCount == 0 {
		fmt.Printf("  %s %s\n\n",
			styleSuccess.Render("✔"),
			styleSuccess.Render(fmt.Sprintf("¡Listo! %d archivos procesados con éxito.", processed)),
		)
	} else {
		fmt.Printf("  %s %s\n\n",
			styleWarning.Render("⚠"),
			styleWarning.Render(fmt.Sprintf("Operación finalizada: %d procesados, %d errores.", processed, errorCount)),
		)
	}
	pause()
}

func showError(err error) {
	fmt.Printf("\n  %s %s\n\n",
		styleError.Render("✖ Error:"),
		styleError.Render(err.Error()),
	)
	pause()
}

func pause() {
	var dummy string
	_ = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Presiona Enter para continuar").
				Options(huh.NewOption("Continuar", "ok")).
				Value(&dummy),
		),
	).WithTheme(CustomTheme()).Run()
}

func init() {
	_ = os.Setenv("TERM", "xterm-256color")
}
