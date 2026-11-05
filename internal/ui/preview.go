package ui

import (
	"fmt"
	"go-archkit/internal/core"

	"github.com/charmbracelet/huh"
)

type PreviewAction string

const (
	ActionApply  PreviewAction = "apply"
	ActionEdit   PreviewAction = "edit"
	ActionCancel PreviewAction = "cancel"
)

// ShowPreview renderiza la vista previa de cambios y solicita confirmación interactiva.
func ShowPreview(plan *core.FilePlan) (PreviewAction, error) {
	fmt.Print(RenderHeader("Vista Previa"))

	total := plan.Total()
	modified := plan.ModifiedCount()

	if total == 0 || modified == 0 {
		fmt.Printf("  %s\n\n", styleSubtitle.Render("No se encontraron archivos que requieran cambios."))
		var dummy string
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Acción").
					Options(huh.NewOption("Regresar al menú principal", "back")).
					Value(&dummy),
			),
		).WithTheme(CustomTheme()).Run()
		return ActionCancel, err
	}

	fmt.Printf("  Se modificarán %d de %d archivos en %s:\n\n",
		modified, total, styleSubtitle.Render(plan.BaseDir))

	// Mostrar muestra de archivos (máximo 20 para no desbordar terminal)
	limit := 20
	shown := 0
	for _, c := range plan.Changes {
		if !c.IsModified() {
			continue
		}
		if shown < limit {
			fmt.Printf("    %s  %s  %s\n",
				styleOldName.Render(c.OldName),
				styleArrow.Render("➜"),
				styleNewName.Render(c.NewName),
			)
			shown++
		}
	}

	if modified > limit {
		fmt.Printf("    %s\n", styleSubtitle.Render(fmt.Sprintf("... y %d archivos más.", modified-limit)))
	}
	fmt.Println()

	// Validar colisiones
	collisions := plan.DetectCollisions()
	if len(collisions) > 0 {
		fmt.Println("  " + styleError.Render("⚠ Se detectaron conflictos que impiden aplicar los cambios:"))
		for _, col := range collisions {
			fmt.Printf("    • %s\n", styleWarning.Render(col))
		}
		fmt.Println()

		var choice string
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("¿Qué deseas hacer?").
					Options(
						huh.NewOption("Volver a editar parámetros", string(ActionEdit)),
						huh.NewOption("Cancelar operación", string(ActionCancel)),
					).
					Value(&choice),
			),
		).WithTheme(CustomTheme()).Run()

		return PreviewAction(choice), err
	}

	// Menú de confirmación sin colisiones
	var action string
	options := []huh.Option[string]{
		huh.NewOption("Sí, aplicar cambios", string(ActionApply)),
		huh.NewOption("Volver a editar parámetros", string(ActionEdit)),
		huh.NewOption("Cancelar", string(ActionCancel)),
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("¿Deseas aplicar estos cambios?").
				Options(options...).
				Value(&action),
		),
	).WithTheme(CustomTheme()).Run()

	if err != nil {
		return ActionCancel, err
	}

	return PreviewAction(action), nil
}
