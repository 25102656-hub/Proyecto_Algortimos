package rtree

import (
	"fmt"
	"strings"
)

func (a *ArbolR) MostrarArbol() {

	fmt.Println("\n===== ESTRUCTURA DEL R-TREE =====")

	a.imprimirNodo(a.Raiz, 0)
}

func (a *ArbolR) imprimirNodo(nodo *Nodo, nivel int) {

	if nodo == nil {
		return
	}

	espacios := strings.Repeat("  ", nivel)

	if nodo.EsHoja {
		fmt.Println(espacios + "Nodo hoja:")
	} else {
		fmt.Println(espacios + "Nodo interno:")
	}

	fmt.Printf(
		"%sZona: [%.4f, %.4f] - [%.4f, %.4f]\n",
		espacios,
		nodo.Zona.MinX,
		nodo.Zona.MinY,
		nodo.Zona.MaxX,
		nodo.Zona.MaxY,
	)

	for _, registro := range nodo.Registros {

		if nodo.EsHoja {

			fmt.Printf(
				"%s  - %s (%s) [ID: %d]\n",
				espacios,
				registro.Lugar.Nombre,
				registro.Lugar.Ciudad,
				registro.Lugar.Id,
			)

		} else {

			a.imprimirNodo(registro.Hijo, nivel+1)
		}
	}
}

func (a *ArbolR) TextoArbol() string {
	var sb strings.Builder
	sb.WriteString("===== ESTRUCTURA DEL R-TREE =====\n")
	a.textoNodo(a.Raiz, 0, &sb)
	return sb.String()
}

func (a *ArbolR) textoNodo(nodo *Nodo, nivel int, sb *strings.Builder) {
	if nodo == nil {
		return
	}

	espacios := strings.Repeat("  ", nivel)

	if nodo.EsHoja {
		sb.WriteString(espacios + "Nodo hoja:\n")
	} else {
		sb.WriteString(espacios + "Nodo interno:\n")
	}

	sb.WriteString(fmt.Sprintf(
		"%sZona: [%.4f, %.4f] - [%.4f, %.4f]\n",
		espacios,
		nodo.Zona.MinX,
		nodo.Zona.MinY,
		nodo.Zona.MaxX,
		nodo.Zona.MaxY,
	))

	for _, registro := range nodo.Registros {
		if nodo.EsHoja {
			sb.WriteString(fmt.Sprintf(
				"%s  - ID: %d | %s | %s | Lat: %.6f | Long: %.6f\n",
				espacios,
				registro.Lugar.Id,
				registro.Lugar.Nombre,
				registro.Lugar.Ciudad,
				registro.Lugar.Latitud,
				registro.Lugar.Longitud,
			))
		} else {
			a.textoNodo(registro.Hijo, nivel+1, sb)
		}
	}
}
