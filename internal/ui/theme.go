package ui

import (
	"fmt"
	"strings"

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
	b.WriteString(fmt.Sprintf("  %s\n\n", styleDivider.Render(strings.Repeat("─", 48))))
	return b.String()
}

// CustomTheme retorna un tema de Huh minimalista y sin recuadros pesados.
func CustomTheme() *huh.Theme {
	t := huh.ThemeBase()

	// Personalización sin bordes
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

	return t
}
