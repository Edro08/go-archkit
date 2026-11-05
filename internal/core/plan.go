package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileChange representa la modificación planificada para un archivo individual.
type FileChange struct {
	SourcePath string // Ruta completa o relativa original
	TargetPath string // Ruta completa o relativa calculada
	OldName    string // Nombre base original
	NewName    string // Nombre base resultante
	Size       int64  // Tamaño del archivo en bytes
	IsDir      bool   // Si es directorio
}

// IsModified indica si el nombre o la ruta cambian respecto al origen.
func (c FileChange) IsModified() bool {
	return filepath.Clean(c.SourcePath) != filepath.Clean(c.TargetPath)
}

// FilePlan almacena el conjunto de operaciones planificadas en memoria antes de ejecutarse.
type FilePlan struct {
	BaseDir string
	Changes []FileChange
}

// NewFilePlan inicializa un nuevo plan de operaciones.
func NewFilePlan(baseDir string) *FilePlan {
	return &FilePlan{
		BaseDir: baseDir,
		Changes: make([]FileChange, 0),
	}
}

// Add añade un cambio al plan.
func (p *FilePlan) Add(change FileChange) {
	p.Changes = append(p.Changes, change)
}

// Total retorna el número total de archivos en el plan.
func (p *FilePlan) Total() int {
	return len(p.Changes)
}

// ModifiedCount cuenta cuántos archivos efectivamente cambiarán.
func (p *FilePlan) ModifiedCount() int {
	count := 0
	for _, c := range p.Changes {
		if c.IsModified() {
			count++
		}
	}
	return count
}

// DetectCollisions verifica si el plan provocaría colisiones de nombres:
// 1. Dos archivos origen diferentes mapeando al mismo destino (colisión interna).
// 2. Un destino que ya existe en el disco y no es el propio archivo origen (colisión externa).
func (p *FilePlan) DetectCollisions() []string {
	var collisions []string
	targetsSeen := make(map[string]string) // targetPath -> sourcePath

	for _, c := range p.Changes {
		if !c.IsModified() {
			continue
		}

		targetClean := filepath.Clean(c.TargetPath)
		sourceClean := filepath.Clean(c.SourcePath)

		// 1. Colisión interna en el plan
		if prevSource, exists := targetsSeen[strings.ToLower(targetClean)]; exists {
			collisions = append(collisions, fmt.Sprintf("Conflicto interno: %q y %q tienen el mismo destino %q",
				filepath.Base(prevSource), filepath.Base(sourceClean), c.NewName))
			continue
		}
		targetsSeen[strings.ToLower(targetClean)] = sourceClean

		// 2. Colisión externa con archivos ya existentes en disco
		if info, err := os.Stat(targetClean); err == nil {
			// Si existe en disco y no es exactamente el mismo archivo origen
			if !os.SameFile(info, getFileInfoOrNil(sourceClean)) {
				collisions = append(collisions, fmt.Sprintf("El archivo destino ya existe en disco: %q", c.NewName))
			}
		}
	}

	return collisions
}

func getFileInfoOrNil(path string) os.FileInfo {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}
	return info
}
