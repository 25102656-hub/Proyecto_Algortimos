package rtree

type ArbolR struct {
	Raiz *Nodo

	Maximo int

	Minimo int
}

func NuevoArbol() *ArbolR {

	return &ArbolR{

		Raiz: NuevoNodo(true),

		Maximo: 4,

		Minimo: 2,
	}

}
