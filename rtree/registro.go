package rtree

import "proyecto_algoritmos/modelos"

type Registro struct {
	Zona Rectangulo

	Lugar modelos.Lugar

	Hijo *Nodo
}
