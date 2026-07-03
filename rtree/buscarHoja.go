package rtree

func (a *ArbolR) buscarHoja(
	nodo *Nodo,
	registro *Registro,
) *Nodo {

	if nodo.EsHoja {
		return nodo
	}

	var mejorRegistro *Registro
	menorCrecimiento := -1.0

	for _, actual := range nodo.Registros {

		crecimiento := actual.Zona.Crecimiento(registro.Zona)

		if mejorRegistro == nil {

			mejorRegistro = actual
			menorCrecimiento = crecimiento
			continue

		}

		if crecimiento < menorCrecimiento {

			mejorRegistro = actual
			menorCrecimiento = crecimiento

		} else if crecimiento == menorCrecimiento {

			if actual.Zona.Area() < mejorRegistro.Zona.Area() {

				mejorRegistro = actual

			}
		}
	}
	return a.buscarHoja(
		mejorRegistro.Hijo,
		registro,
	)

}
