package rtree

func (a *ArbolR) EliminarPorID(id int) bool {
	return a.eliminarEnNodo(a.Raiz, id)
}

func (a *ArbolR) eliminarEnNodo(nodo *Nodo, id int) bool {
	if nodo == nil {
		return false
	}

	if nodo.EsHoja {
		for i, registro := range nodo.Registros {
			if registro.Lugar.Id == id {
				nodo.EliminarRegistroPorIndice(i)
				a.actualizarHaciaArriba(nodo)
				return true
			}
		}

		return false
	}

	for _, registro := range nodo.Registros {
		encontrado := a.eliminarEnNodo(registro.Hijo, id)

		if encontrado {
			if registro.Hijo != nil && len(registro.Hijo.Registros) == 0 {
				a.quitarHijoVacio(nodo, registro.Hijo)
			}

			a.actualizarHaciaArriba(nodo)
			return true
		}
	}

	return false
}

func (a *ArbolR) quitarHijoVacio(padre *Nodo, hijo *Nodo) {
	for i, registro := range padre.Registros {
		if registro.Hijo == hijo {
			padre.EliminarRegistroPorIndice(i)
			return
		}
	}
}

func (a *ArbolR) actualizarHaciaArriba(nodo *Nodo) {
	for nodo != nil {
		nodo.ActualizarZona()
		a.actualizarRegistroEnPadre(nodo)
		nodo = nodo.Padre
	}
}
