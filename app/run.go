package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Run() {
	reader := bufio.NewReader(os.Stdin)

	for {
		showMenu()

		option, _ := reader.ReadString('\n')
		switch strings.TrimSpace(option) {
		case "1":
			if err := NewChangeExt(reader).Change(); err != nil {
				fmt.Println("Error:", err)
			}
		case "2":
			if err := NewFileNumber(reader).Number(); err != nil {
				fmt.Println("Error:", err)
			}
		case "3":
			if err := NewFileRename(reader).Rename(); err != nil {
				fmt.Println("Error:", err)
			}
		case "0":
			fmt.Println("Saliendo...")
			return

		default:
			fmt.Println("Opción no válida")
		}

		fmt.Println()
	}
}

func showMenu() {
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Println("Archkit v1.0.0 by Edro08")
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Println("Opciones:")
	fmt.Println("  1. Cambiar extensiones")
	fmt.Println("  2. Numerar archivos")
	fmt.Println("  3. Renombrar archivos")
	fmt.Println("  0. Salir")
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Print("Seleccione una opción: ")
}
