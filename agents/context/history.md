# Historial de Archkit

## 2026-11-04 — Fases 3 y 4: TUI Minimalista con Charmbracelet, Vista Previa Diff y Retiro de `app/`

- **Retiro y Erradicación de `app/`**:
  - Eliminado por completo el paquete legado `app/` (`files_ext_change.go`, `files_num.go`, `files_rename.go`, `run.go`, `utils.go`) eliminando el acoplamiento y código obsoleto.
- **Tema y Estética Minimalista ([`internal/ui/theme.go`](file:///d:/Projects/Go/go-archkit/internal/ui/theme.go))**:
  - Implementación de `RenderHeader` uniforme `Archkit v2.0 › [Sección]` con separador de texto sutil.
  - Configuración de `CustomTheme()` sobre `huh.ThemeBase()` sin marcos pesados, con sangrías suaves y paleta de acentos Índigo y Cyan.
- **Formularios Guiados e Inteligentes ([`internal/ui/forms.go`](file:///d:/Projects/Go/go-archkit/internal/ui/forms.go))**:
  - Formulario de extensiones con escaneo automático en vivo de extensiones presentes en la carpeta seleccionada (`ScanExtensions`), permitiendo elegir con flechas del teclado en lugar de escribir a ciegas.
  - Formularios dedicados para numeración secuencial, búsqueda/reemplazo y normalización de formatos (`kebab-case`, `snake_case`, etc.) con validación de directorios y números enteros.
- **Vista Previa Diff Interactiva ([`internal/ui/preview.go`](file:///d:/Projects/Go/go-archkit/internal/ui/preview.go))**:
  - Renderizado de diferencias estilo Git: nombre original en atenuado, flecha `➜` y nombre destino en verde destacado.
  - Alerta inmediata de colisiones de nombres impidiendo la ejecución en caso de conflicto.
  - Menú de confirmación con opción de retorno («Volver a editar») sin perder los valores ingresados.
- **Orquestador Principal y Entrypoint ([`internal/ui/app.go`](file:///d:/Projects/Go/go-archkit/internal/ui/app.go), [`main.go`](file:///d:/Projects/Go/go-archkit/main.go))**:
  - Reemplazo del bucle principal en `main.go` delegando en `ui.Run()`.
- **Refinamiento de UX, Barra de Atajos y Valores por Defecto ([`internal/ui/theme.go`](file:///d:/Projects/Go/go-archkit/internal/ui/theme.go), [`internal/ui/forms.go`](file:///d:/Projects/Go/go-archkit/internal/ui/forms.go), [`README.md`](file:///d:/Projects/Go/go-archkit/README.md))**:
  - Incorporada barra de ayuda inferior unificada (`FooterNote`) presente en todas las pantallas (`↑/↓ Navigate • ←/→ Toggle • Enter Submit • Esc Back`).
  - Habilitada cancelación instantánea con tecla `Esc` en cualquier campo sin necesidad de completar el formulario.
  - Ajustados valores predeterminados seguros en numeración: conservar nombre original en `No`, in-place en `No` y subcarpeta destino por defecto en `temp/`.
  - Limpieza de pantalla y scrollback (`ClearScreen`) entre vistas para evitar encabezados duplicados.
  - Unificada la sección de compilación y ejecución rápida en `README.md`.
- **Verificación**:
  - Compilación exitosa del binario `archkit.exe` (`go build -o archkit.exe .`) y ejecución de pruebas unitarias (`go test ./...`) aprobadas.

## 2026-11-04 — Fase 2: Arquitectura Desacoplada y Motor de Archivos Core

- **Modelos de Planificación y Prevención de Colisiones ([`internal/core/plan.go`](file:///d:/Projects/Go/go-archkit/internal/core/plan.go))**:
  - Creación de estructuras `FileChange` y `FilePlan` para representar las mutaciones de archivos calculadas puramente en memoria.
  - Método `DetectCollisions()` que valida colisiones internas (dos orígenes con idéntico destino) y externas (conflictos con archivos preexistentes en disco).
- **Ordenación Natural ([`internal/core/natural_sort.go`](file:///d:/Projects/Go/go-archkit/internal/core/natural_sort.go))**:
  - Implementación de `NaturalLess` desacoplada para garantizar orden numérico humano (`2.txt` precede a `10.txt`).
- **Procesadores Core de Archivos ([`internal/core/`](file:///d:/Projects/Go/go-archkit/internal/core/))**:
  - [`ext_changer.go`](file:///d:/Projects/Go/go-archkit/internal/core/ext_changer.go): `PlanExtChange` y escaneo previo con `ScanExtensions` para estadísticas de extensiones existentes.
  - [`numberer.go`](file:///d:/Projects/Go/go-archkit/internal/core/numberer.go): `PlanNumbering` con soporte de padding dinámico, prefijos, sufijos y opción in-place o subcarpeta.
  - [`renamer.go`](file:///d:/Projects/Go/go-archkit/internal/core/renamer.go): `PlanRename` con búsqueda/reemplazo (case-sensitive o insensitive) y decoradores.
  - [`normalizer.go`](file:///d:/Projects/Go/go-archkit/internal/core/normalizer.go): `PlanNormalize` y `ConvertCase` para kebab-case, snake_case, camelCase y PascalCase.
- **Ejecutor Seguro ([`internal/core/executor.go`](file:///d:/Projects/Go/go-archkit/internal/core/executor.go))**:
  - Aplicación física de planes en disco mediante `Execute` con modos `rename` o `copy`, creación automática de directorios y control de errores atómico.
- **Verificación**:
  - Suite de pruebas unitarias completa en [`internal/core/core_test.go`](file:///d:/Projects/Go/go-archkit/internal/core/core_test.go) (`go test ./internal/core/...`) ejecutada exitosamente con 100% de tests aprobados.

## 2026-10-16 — Fase 1: Creación y Funcionalidades Base v1.0.0 (Asistente de Consola)

- **Menú y Flujo Interactivo Base ([`app/run.go`](file:///d:/Projects/Go/go-archkit/app/run.go), [`app/utils.go`](file:///d:/Projects/Go/go-archkit/app/utils.go))**:
  - Implementación de menú de consola interactivo secuencial mediante `bufio.Reader` leyendo de `os.Stdin`.
  - Rutinas auxiliares para selección de directorio de trabajo, validación de extensiones y diálogo de confirmación `(S/N)`.
- **Cambio Masivo de Extensiones ([`app/files_ext_change.go`](file:///d:/Projects/Go/go-archkit/app/files_ext_change.go))**:
  - Detección de archivos regulares en directorio plano coincidentes con extensión origen.
  - Renombrado *in-place* con cálculo de nueva extensión y confirmación previa de cambios.
- **Numeración Secuencial de Archivos ([`app/files_num.go`](file:///d:/Projects/Go/go-archkit/app/files_num.go))**:
  - Algoritmo de ordenación natural (`naturalLess`) para secuencias numéricas en nombres de archivo.
  - Soporte de números iniciales arbitrarios con `big.Int` y padding dinámico de ceros a la izquierda (mínimo 3 dígitos).
  - Aislamiento seguro: copia los archivos numerados a una subcarpeta `temp/` para evitar colisiones destructivas.
- **Búsqueda, Reemplazo y Prefijos/Sufijos ([`app/files_rename.go`](file:///d:/Projects/Go/go-archkit/app/files_rename.go))**:
  - Opciones combinadas para filtrar por extensión, buscar y reemplazar cadenas de texto, y añadir prefijos y sufijos en nombres de archivo.
- **Compilación**:
  - Generación del ejecutable `archkit.exe` en Go puro sin dependencias externas.
