package rtree

type Nodo struct {
	EsHoja bool
	Padre  *Nodo

	Registros []*Registro
	Zona      Rectangulo
}

func NuevoNodo(esHoja bool) *Nodo {
	return &Nodo{
		EsHoja:    esHoja,
		Registros: []*Registro{},
	}
}

func (n *Nodo) ActualizarZona() {
	if len(n.Registros) == 0 {
		n.Zona = Rectangulo{}
		return
	}

	zona := n.Registros[0].Zona

	for i := 1; i < len(n.Registros); i++ {
		zona = zona.Unir(n.Registros[i].Zona)
	}

	n.Zona = zona
}

func (n *Nodo) EsRaiz() bool {
	return n.Padre == nil
}

func (n *Nodo) CantidadRegistros() int {
	return len(n.Registros)
}

func (n *Nodo) AgregarRegistro(registro *Registro) {
	n.Registros = append(n.Registros, registro)

	if registro.Hijo != nil {
		registro.Hijo.Padre = n
	}

	n.ActualizarZona()
}

func (n *Nodo) EliminarRegistroPorIndice(indice int) {
	n.Registros = append(n.Registros[:indice], n.Registros[indice+1:]...)
	n.ActualizarZona()
}
