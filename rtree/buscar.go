package rtree

import "proyecto_algoritmos/modelos"

func (a *ArbolR) BuscarPorRango(zonaBusqueda Rectangulo) []modelos.Lugar {

	var resultados []modelos.Lugar

	a.buscarEnNodo(
		a.Raiz,
		zonaBusqueda,
		&resultados,
	)

	return resultados
}

func (a *ArbolR) buscarEnNodo(
	nodo *Nodo,
	zonaBusqueda Rectangulo,
	resultados *[]modelos.Lugar,
) {

	if nodo == nil {
		return
	}

	for _, registro := range nodo.Registros {

		if !registro.Zona.Intersecta(zonaBusqueda) {
			continue
		}

		if nodo.EsHoja {

			*resultados = append(
				*resultados,
				registro.Lugar,
			)

		} else {

			a.buscarEnNodo(
				registro.Hijo,
				zonaBusqueda,
				resultados,
			)

		}
	}
}
