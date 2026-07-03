package rtree

import (
	"testing"

	"proyecto_algoritmos/modelos"
)

func TestInsertarYBuscarPorRango(t *testing.T) {
	arbol := NuevoArbol()

	lugar := modelos.Lugar{
		Id:        1,
		Nombre:    "Plaza Mayor",
		Categoria: "Histórico",
		Ciudad:    "Lima",
		Latitud:   -12.046374,
		Longitud:  -77.042793,
	}

	arbol.Insertar(lugar)

	zona := Rectangulo{
		MinX: -77.10,
		MinY: -12.10,
		MaxX: -77.00,
		MaxY: -12.00,
	}

	resultados := arbol.BuscarPorRango(zona)

	if len(resultados) == 0 {
		t.Fatalf("Se esperaba encontrar el lugar insertado, pero no se encontró.")
	}

	if resultados[0].Id != lugar.Id {
		t.Errorf("Se esperaba ID %d, pero se obtuvo ID %d", lugar.Id, resultados[0].Id)
	}
}

func TestEliminarPorID(t *testing.T) {
	arbol := NuevoArbol()

	lugar := modelos.Lugar{
		Id:        2,
		Nombre:    "Parque Kennedy",
		Categoria: "Parque",
		Ciudad:    "Lima",
		Latitud:   -12.121120,
		Longitud:  -77.029700,
	}

	arbol.Insertar(lugar)

	eliminado := arbol.EliminarPorID(2)

	if !eliminado {
		t.Errorf("Se esperaba eliminar el lugar con ID 2.")
	}

	zona := Rectangulo{
		MinX: -77.10,
		MinY: -12.20,
		MaxX: -77.00,
		MaxY: -12.00,
	}

	resultados := arbol.BuscarPorRango(zona)

	if len(resultados) != 0 {
		t.Errorf("El lugar fue eliminado, pero todavía aparece en la búsqueda.")
	}
}

func BenchmarkInsertar1000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		arbol := NuevoArbol()

		for i := 0; i < 1000; i++ {
			lugar := modelos.Lugar{
				Id:        i,
				Nombre:    "Lugar Demo",
				Categoria: "Prueba",
				Ciudad:    "Lima",
				Latitud:   -12.0 + float64(i)*0.000001,
				Longitud:  -77.0 + float64(i)*0.000001,
			}

			arbol.Insertar(lugar)
		}
	}
}

func BenchmarkBuscarPorRango1000(b *testing.B) {
	arbol := NuevoArbol()

	for i := 0; i < 1000; i++ {
		lugar := modelos.Lugar{
			Id:        i,
			Nombre:    "Lugar Demo",
			Categoria: "Prueba",
			Ciudad:    "Lima",
			Latitud:   -12.0 + float64(i)*0.000001,
			Longitud:  -77.0 + float64(i)*0.000001,
		}

		arbol.Insertar(lugar)
	}

	zona := Rectangulo{
		MinX: -77.10,
		MinY: -12.10,
		MaxX: -76.90,
		MaxY: -11.90,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		arbol.BuscarPorRango(zona)
	}
}
