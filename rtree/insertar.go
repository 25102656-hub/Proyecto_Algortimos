package rtree

import "proyecto_algoritmos/modelos"

func (a *ArbolR) Insertar(lugar modelos.Lugar) {

	zona := NuevoRectangulo(

		lugar.Longitud,

		lugar.Latitud,
	)

	registro := &Registro{

		Zona: zona,

		Lugar: lugar,
	}

	hoja := a.buscarHoja(

		a.Raiz,

		registro,
	)

	hoja.Registros = append(

		hoja.Registros,

		registro,
	)

	hoja.ActualizarZona()

	a.ajustarPadres(

		hoja,
	)

}
