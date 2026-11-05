# Archkit

Herramienta de consola para Windows y multiplataforma escrita en Go, diseñada para operaciones por lote seguras y rápidas sobre archivos.

Cuenta con una **TUI minimalista e interactiva** (construida con [Charm](https://charm.sh): `huh` y `lipgloss`) que te guía paso a paso sin necesidad de memorizar argumentos de línea de comandos.

---

## Características

* **Cambiar extensiones:** Escaneo automático de extensiones en la carpeta elegida, selección con flechas y soporte recursivo opcional.
* **Numerar en serie:** Ordenamiento natural (`item_2.txt` precede a `item_10.txt`), padding de ceros a la izquierda automático y opción de renombrar *in-place* o mover a una subcarpeta.
* **Buscar y reemplazar:** Reemplazo de cadenas (con opción sensible o insensible a mayúsculas) y adición de prefijos y sufijos.
* **Normalizar nombres:** Conversión masiva de formatos a `kebab-case`, `snake_case`, `camelCase`, `PascalCase`, minúsculas o mayúsculas.
* **Vista previa segura (Diff):** Previsualiza exactamente qué archivos cambiarán antes de tocar el disco, con detección anticipada de colisiones o conflictos.

---

## Compilación

Requiere Go 1.25 o posterior:

```powershell
go build -o archkit.exe .
```

## Uso

Simplemente ejecuta el binario en la terminal o haz doble clic:

```powershell
.\archkit.exe
```
