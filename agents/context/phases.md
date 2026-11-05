# Fases de desarrollo de Archkit

La fase 1 representa la versión original de consola interactiva basada en menús secuenciales de texto plano. Las fases 2 a 5 definen la hoja de ruta técnica para la versión 2.0 basada en una TUI minimalista y fluida con Charmbracelet (`huh` y `lipgloss`).

---

## 🟢 Fases completadas

### Fase 1 — Motor base y CLI interactiva v1.0.0
- Asistente de terminal interactivo secuencial mediante `bufio.Reader`.
- Cambio de extensiones masivo en directorio local.
- Numeración de archivos con ordenamiento natural (`naturalLess`) y cálculo de padding dinámico en carpeta `temp/`.
- Renombrado por búsqueda, reemplazo, prefijos y sufijos.
- Cero dependencias externas en Go puro.

### Fase 2 — Arquitectura desacoplada y Motor de Archivos Core
- Motor de dominio puro desacoplado en `internal/core` con cero dependencias externas.
- Generación de planes en memoria (`FilePlan` / `FileChange`) que calculan el diff antes de tocar el disco.
- Detección preventiva de colisiones de nombres (internas y externas).
- Algoritmo de ordenación natural `NaturalLess` independiente y reutilizable.
- Procesadores de extensión (`PlanExtChange`, `ScanExtensions`), numeración (`PlanNumbering`), renombrado (`PlanRename`) y normalización de formatos (`PlanNormalize` con kebab, snake, camel, pascal).
- Ejecutor resiliente (`Execute`) con soporte de modos rename *in-place* y copia a directorio alternativo.
- Cobertura de pruebas unitarias al 100% en `internal/core/core_test.go`.

### Fase 3 y 4 — TUI Minimalista con Charmbracelet y Preview Diff Interactivo
- Erradicación y eliminación definitiva del paquete legado `app/`.
- Nueva capa de presentación `internal/ui/` con tema minimalista unificado sin cajas ni recuadros pesados (`theme.go`).
- Formularios interactivos guiados con `charmbracelet/huh` para todas las operaciones (`forms.go`), incluyendo autocompletado y detección dinámica de extensiones existentes.
- Pantalla de Vista Previa (*Diff*) estilo Git con flechas `➜`, formato de colores semánticos (`lipgloss`), alerta de colisiones y ciclo de retorno («Volver a editar») sin perder datos (`preview.go`).
- Orquestador interactivo en `app.go` conectado con el nuevo `main.go`.

### Fase 5 — Pruebas Finales, Verificación y Documentación
- Verificación funcional integral con pruebas manuales en escenarios de prueba (extensiones, ordenamiento natural, normalización y detección de colisiones).
- Compilación del binario `archkit.exe`.
- Actualización de documentación en `README.md`.
