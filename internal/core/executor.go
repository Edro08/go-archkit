package core

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileExecutionError registra un fallo en un archivo particular.
type FileExecutionError struct {
	File string
	Err  error
}

func (e FileExecutionError) Error() string {
	return fmt.Sprintf("%s: %v", e.File, e.Err)
}

// ExecutionResult resume el resultado de ejecutar un plan.
type ExecutionResult struct {
	Processed int
	Skipped   int
	Errors    []FileExecutionError
}

// ExecuteMode indica si la ejecución mueve (renombra) o copia los archivos.
type ExecuteMode string

const (
	ModeRename ExecuteMode = "rename" // Renombra o mueve el archivo
	ModeCopy   ExecuteMode = "copy"   // Copia el archivo preservando el original
)

// Execute aplica un FilePlan en el disco.
func Execute(plan *FilePlan, mode ExecuteMode) (*ExecutionResult, error) {
	if plan == nil {
		return nil, errors.New("plan nulo")
	}

	// 1. Validar colisiones antes de tocar el sistema de archivos
	collisions := plan.DetectCollisions()
	if len(collisions) > 0 {
		return nil, fmt.Errorf("se detectaron conflictos que impiden la ejecución: %s", collisions[0])
	}

	result := &ExecutionResult{
		Errors: make([]FileExecutionError, 0),
	}

	for _, change := range plan.Changes {
		if !change.IsModified() {
			result.Skipped++
			continue
		}

		// Asegurar que el directorio de destino exista
		targetDir := filepath.Dir(change.TargetPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			result.Errors = append(result.Errors, FileExecutionError{
				File: change.OldName,
				Err:  fmt.Errorf("no se pudo crear directorio destino: %w", err),
			})
			continue
		}

		var opErr error
		if mode == ModeCopy {
			opErr = copyFile(change.SourcePath, change.TargetPath)
		} else {
			opErr = os.Rename(change.SourcePath, change.TargetPath)
		}

		if opErr != nil {
			result.Errors = append(result.Errors, FileExecutionError{
				File: change.OldName,
				Err:  opErr,
			})
		} else {
			result.Processed++
		}
	}

	return result, nil
}

func copyFile(source, destination string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if closeErr := out.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(destination)
	}
	return err
}
