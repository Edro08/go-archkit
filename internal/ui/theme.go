package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Paleta de colores minimalista
	colorDim     = lipgloss.AdaptiveColor{Light: "#94A3B8", Dark: "#64748B"}
	colorPrimary = lipgloss.AdaptiveColor{Light: "#4F46E5", Dark: "#818CF8"} // Índigo
	colorAccent  = lipgloss.AdaptiveColor{Light: "#0284C7", Dark: "#38BDF8"} // Cyan
	colorSuccess = lipgloss.AdaptiveColor{Light: "#16A34A", Dark: "#4ADE80"} // Verde
	colorWarning = lipgloss.AdaptiveColor{Light: "#EA580C", Dark: "#FB923C"} // Naranja
	colorError   = lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#F87171"} // Rojo
	colorText    = lipgloss.AdaptiveColor{Light: "#1E293B", Dark: "#F1F5F9"}

	// Estilos tipográficos
	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary)

	styleSubtitle = lipgloss.NewStyle().
			Foreground(colorDim)

	styleDivider = lipgloss.NewStyle().
			Foreground(colorDim)

	styleArrow = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	styleOldName = lipgloss.NewStyle().
			Foreground(colorDim)

	styleNewName = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Bold(true)

	styleSuccess = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Bold(true)

	styleError = lipgloss.NewStyle().
			Foreground(colorError).
			Bold(true)

	styleWarning = lipgloss.NewStyle().
			Foreground(colorWarning).
			Bold(true)
)

// ClearScreen limpia la terminal y el scrollback para evitar acumulaciones de cabeceras.
func ClearScreen() {
	fmt.Print("\033[H\033[2J\033[3J")
}

// RenderHeader imprime el encabezado minimalista unificado con la sección activa.
func RenderHeader(section string) string {
	var b strings.Builder
	b.WriteString("\n")
	if section != "" {
		b.WriteString(fmt.Sprintf("  %s %s %s\n",
			styleTitle.Render("Archkit v2.0"),
			styleSubtitle.Render("›"),
			lipgloss.NewStyle().Bold(true).Foreground(colorText).Render(section),
		))
	} else {
		b.WriteString(fmt.Sprintf("  %s\n", styleTitle.Render("Archkit v2.0")))
	}
	b.WriteString(fmt.Sprintf("  %s\n\n", styleDivider.Render(strings.Repeat("─", 54))))
	return b.String()
}

// FooterNote crea el componente visual de ayuda unificado para el pie de todas las pantallas.
func FooterNote() *huh.Note {
	divider := styleDivider.Render(strings.Repeat("─", 54))
	keys := fmt.Sprintf("%s %s  •  %s %s  •  %s %s  •  %s %s",
		lipgloss.NewStyle().Bold(true).Foreground(colorPrimary).Render("↑/↓"),
		styleSubtitle.Render("Navigate"),
		lipgloss.NewStyle().Bold(true).Foreground(colorPrimary).Render("←/→"),
		styleSubtitle.Render("Toggle"),
		lipgloss.NewStyle().Bold(true).Foreground(colorPrimary).Render("Enter"),
		styleSubtitle.Render("Submit"),
		lipgloss.NewStyle().Bold(true).Foreground(colorPrimary).Render("Esc"),
		styleSubtitle.Render("Back"),
	)
	footer := fmt.Sprintf("\n%s\n  %s", divider, keys)
	return huh.NewNote().Title(footer)
}

// CustomTheme retorna un tema de Huh minimalista, compacto y sin dobles saltos de línea.
func CustomTheme() *huh.Theme {
	t := huh.ThemeBase()

	// Separador compacto: 1 solo salto de línea en lugar del doble salto (\n\n) por defecto
	t.FieldSeparator = lipgloss.NewStyle().SetString("\n")

	// Alinear campos a 2 espacios
	t.Focused.Base = lipgloss.NewStyle().PaddingLeft(2)
	t.Blurred.Base = lipgloss.NewStyle().PaddingLeft(2)

	t.Focused.Title = lipgloss.NewStyle().Bold(true).Foreground(colorPrimary)
	t.Blurred.Title = lipgloss.NewStyle().Foreground(colorDim)

	t.Focused.Description = lipgloss.NewStyle().Foreground(colorDim)
	t.Blurred.Description = lipgloss.NewStyle().Foreground(colorDim)

	// Opciones seleccionadas
	t.Focused.SelectedOption = lipgloss.NewStyle().Foreground(colorSuccess).Bold(true)
	t.Focused.UnselectedOption = lipgloss.NewStyle().Foreground(colorText)

	// Inputs
	t.Focused.TextInput.Prompt = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	t.Focused.TextInput.Text = lipgloss.NewStyle().Foreground(colorText)
	t.Focused.TextInput.Placeholder = lipgloss.NewStyle().Foreground(colorDim)

	// Botones Confirm
	t.Focused.FocusedButton = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(colorPrimary).
		Padding(0, 1).
		MarginRight(1)

	t.Focused.BlurredButton = lipgloss.NewStyle().
		Foreground(colorDim).
		Padding(0, 1).
		MarginRight(1)

	t.Blurred.FocusedButton = t.Focused.FocusedButton
	t.Blurred.BlurredButton = t.Focused.BlurredButton

	// Notas
	t.Focused.NoteTitle = lipgloss.NewStyle()
	t.Blurred.NoteTitle = lipgloss.NewStyle()

	return t
}

// CustomKeyMap configura atajos globales permitiendo cancelar con Esc o Ctrl+C.
func CustomKeyMap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Quit = key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("esc", "cancelar"),
	)
	return km
}
