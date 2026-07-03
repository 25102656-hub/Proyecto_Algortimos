package rtree

func (a *ArbolR) dividirRaiz() {
	nodo1, nodo2 := a.splitCuadratico(a.Raiz)

	nuevaRaiz := NuevoNodo(false)
	nodo1.Padre = nuevaRaiz
	nodo2.Padre = nuevaRaiz

	nuevaRaiz.AgregarRegistro(&Registro{Zona: nodo1.Zona, Hijo: nodo1})
	nuevaRaiz.AgregarRegistro(&Registro{Zona: nodo2.Zona, Hijo: nodo2})
	nuevaRaiz.ActualizarZona()

	a.Raiz = nuevaRaiz
}

func (a *ArbolR) dividirNodo(nodo *Nodo) *Nodo {
	padre := nodo.Padre
	nodo1, nodo2 := a.splitCuadratico(nodo)

	nodo1.Padre = padre
	nodo2.Padre = padre

	for i, registro := range padre.Registros {
		if registro.Hijo == nodo {
			padre.EliminarRegistroPorIndice(i)
			break
		}
	}

	padre.AgregarRegistro(&Registro{Zona: nodo1.Zona, Hijo: nodo1})
	padre.AgregarRegistro(&Registro{Zona: nodo2.Zona, Hijo: nodo2})
	padre.ActualizarZona()

	return padre
}

func (a *ArbolR) splitCuadratico(nodo *Nodo) (*Nodo, *Nodo) {
	semilla1, semilla2 := elegirSemillas(nodo.Registros)

	nodo1 := NuevoNodo(nodo.EsHoja)
	nodo2 := NuevoNodo(nodo.EsHoja)

	nodo1.AgregarRegistro(semilla1)
	nodo2.AgregarRegistro(semilla2)

	for _, registro := range nodo.Registros {
		if registro == semilla1 || registro == semilla2 {
			continue
		}

		a.repartirRegistro(registro, nodo1, nodo2)
	}

	nodo1.ActualizarZona()
	nodo2.ActualizarZona()

	return nodo1, nodo2
}

func elegirSemillas(registros []*Registro) (*Registro, *Registro) {
	var semilla1 *Registro
	var semilla2 *Registro
	mayorDesperdicio := -1.0

	for i := 0; i < len(registros)-1; i++ {
		for j := i + 1; j < len(registros); j++ {
			r1 := registros[i].Zona
			r2 := registros[j].Zona
			union := r1.Unir(r2)

			desperdicio := union.Area() - r1.Area() - r2.Area()

			if desperdicio > mayorDesperdicio {
				mayorDesperdicio = desperdicio
				semilla1 = registros[i]
				semilla2 = registros[j]
			}
		}
	}

	return semilla1, semilla2
}

func (a *ArbolR) repartirRegistro(registro *Registro, nodo1 *Nodo, nodo2 *Nodo) {
	crece1 := nodo1.Zona.Crecimiento(registro.Zona)
	crece2 := nodo2.Zona.Crecimiento(registro.Zona)

	if crece1 < crece2 {
		nodo1.AgregarRegistro(registro)
		return
	}

	if crece2 < crece1 {
		nodo2.AgregarRegistro(registro)
		return
	}

	if nodo1.Zona.Area() < nodo2.Zona.Area() {
		nodo1.AgregarRegistro(registro)
		return
	}

	if nodo2.Zona.Area() < nodo1.Zona.Area() {
		nodo2.AgregarRegistro(registro)
		return
	}

	if len(nodo1.Registros) <= len(nodo2.Registros) {
		nodo1.AgregarRegistro(registro)
	} else {
		nodo2.AgregarRegistro(registro)
	}
}

func (a *ArbolR) ajustarPadres(nodo *Nodo) {
	for nodo != nil {
		nodo.ActualizarZona()
		a.actualizarRegistroEnPadre(nodo)

		if len(nodo.Registros) > a.Maximo {
			if nodo == a.Raiz {
				a.dividirRaiz()
				return
			}

			nodo = a.dividirNodo(nodo)
			continue
		}

		nodo = nodo.Padre
	}
}

func (a *ArbolR) actualizarRegistroEnPadre(nodo *Nodo) {
	if nodo == nil || nodo.Padre == nil {
		return
	}

	for _, registro := range nodo.Padre.Registros {
		if registro.Hijo == nodo {
			registro.Zona = nodo.Zona
			return
		}
	}
}
