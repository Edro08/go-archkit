# Arquitectura de Archkit

## 1. Principios de Diseño
1. **Desacoplamiento Estricto (Core vs TUI)**: El motor de manipulación de archivos (`internal/core`) no debe importar librerías de interfaz de usuario (`huh`, `lipgloss`, etc.). Debe ser 100% testeable mediante pruebas unitarias en memoria y carpetas temporales (`t.TempDir()`).
2. **Planificación Previa (Dry-Run por Diseño)**: Toda operación genera primero un plan de cambios (`FilePlan` / `[]FileChange`) en memoria antes de tocar el disco. Esto permite renderizar la vista previa (diff) y detectar colisiones de nombres o errores antes de mutar el sistema de archivos.
3. **Dirección de Dependencias**:
   - `main.go` / `cmd/` -> orquesta y arranca la aplicación.
   - `internal/ui` (o `internal/tui`) -> consume `internal/core` para presentar formularios, previews y confirmaciones.
   - `internal/core` -> depende únicamente de la biblioteca estándar de Go (`os`, `path/filepath`, `strings`, `sort`, `unicode`).

## 2. Estructura de Paquetes Prevista (v2.0)
```text
go-archkit/
├── main.go                     # Punto de entrada
├── internal/
│   ├── core/                   # Motor de dominio (cero dependencias externas)
│   │   ├── plan.go             # Modelos de cambio (FileChange, FilePlan, colisiones)
│   │   ├── ext_changer.go      # Lógica de cambio de extensiones
│   │   ├── numberer.go         # Lógica de numeración secuencial y natural sorting
│   │   ├── renamer.go          # Lógica de reemplazo, prefijos y sufijos
│   │   ├── normalizer.go       # Lógica de case conversion (kebab, snake, camel)
│   │   └── executor.go         # Aplicación segura de planes en disco
│   └── ui/                     # Capa de presentación (Charmbracelet)
│       ├── app.go              # Orquestador del flujo interactivo
│       ├── theme.go            # Estilos minimalistas unificados (lipgloss)
│       ├── forms/              # Formularios interactivos (huh)
│       └── views/              # Vista previa / diff de cambios
└── agents/                     # Contexto, especificaciones e historial del agente
    ├── context/
    │   ├── history.md
    │   └── fases.md
    └── specs/
        └── architecture.md
```

## 3. Reglas de Validación y Seguridad en Filesystem
- **Rutas canónicas y relativas seguras**: Evitar path traversal no deseado; resolver rutas con `filepath.Abs` y `filepath.Clean`.
- **Detección de Colisiones**: Si dos archivos diferentes terminan con el mismo nombre de destino en el mismo directorio, el plan debe marcar un error de colisión antes de solicitar confirmación.
- **Atomicidad o Reporte Detallado**: Si una operación falla a mitad del proceso, reportar claramente los archivos afectados, los omitidos y el error exacto.
