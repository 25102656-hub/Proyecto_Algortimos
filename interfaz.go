//go:build fyne
// +build fyne

package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"proyecto_algoritmos/modelos"
	"proyecto_algoritmos/rtree"
)

func main() {
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error en Ping SQL Server: ", err)
	}

	arbol := rtree.NuevoArbol()

	lugares, err := cargarLugares(db)
	if err != nil {
		log.Fatal("Error al cargar lugares: ", err)
	}

	for _, lugar := range lugares {
		arbol.Insertar(lugar)
	}

	myApp := app.New()
	ventana := myApp.NewWindow("Proyecto Algoritmos - R-Tree")
	ventana.Resize(fyne.NewSize(950, 700))

	resultado := widget.NewMultiLineEntry()
	resultado.SetText(fmt.Sprintf("Conexión exitosa con SQL Server.\nSe cargaron %d lugares desde la base de datos.", len(lugares)))
	resultado.Wrapping = fyne.TextWrapWord
	resultado.SetMinRowsVisible(16)

	nombreEntry := widget.NewEntry()
	nombreEntry.SetPlaceHolder("Ejemplo: Plaza San Martín")
	categoriaEntry := widget.NewEntry()
	categoriaEntry.SetPlaceHolder("Ejemplo: Lima, Cusco, Arequipa")
	ciudadEntry := widget.NewEntry()
	ciudadEntry.SetPlaceHolder("Ejemplo: Lima")
	latitudEntry := widget.NewEntry()
	latitudEntry.SetPlaceHolder("Ejemplo: -12.051944")
	longitudEntry := widget.NewEntry()
	longitudEntry.SetPlaceHolder("Ejemplo: -77.034722")

	formInsertar := widget.NewForm(
		widget.NewFormItem("Nombre", nombreEntry),
		widget.NewFormItem("Provincia / Departamento", categoriaEntry),
		widget.NewFormItem("Ciudad", ciudadEntry),
		widget.NewFormItem("Latitud", latitudEntry),
		widget.NewFormItem("Longitud", longitudEntry),
	)

	btnInsertar := widget.NewButton("Insertar lugar", func() {
		nombre := strings.TrimSpace(nombreEntry.Text)
		provincia := strings.TrimSpace(categoriaEntry.Text)
		ciudad := strings.TrimSpace(ciudadEntry.Text)
		latitud, ok := leerFloatGUI(latitudEntry.Text, "latitud", -90, 90, resultado)
		if !ok {
			return
		}
		longitud, ok := leerFloatGUI(longitudEntry.Text, "longitud", -180, 180, resultado)
		if !ok {
			return
		}

		if nombre == "" || provincia == "" || ciudad == "" {
			resultado.SetText("Nombre, categoría y ciudad no pueden estar vacíos.")
			return
		}

		var idGenerado int
		err := db.QueryRow(
			"INSERT INTO lugares (Nombre, Categoria, Ciudad, Latitud, Longitud) OUTPUT INSERTED.ID VALUES (@p1, @p2, @p3, @p4, @p5)",
			nombre,
			provincia,
			ciudad,
			latitud,
			longitud,
		).Scan(&idGenerado)

		if err != nil {
			resultado.SetText("Error al insertar en la base de datos: " + err.Error() + "\nLa tabla debe tener la columna Categoria.")
			return
		}

		lugar := modelos.Lugar{
			Id:        idGenerado,
			Nombre:    nombre,
			Categoria: provincia,
			Ciudad:    ciudad,
			Latitud:   latitud,
			Longitud:  longitud,
		}

		arbol.Insertar(lugar)
		resultado.SetText(fmt.Sprintf("insertado correctamente.\nID: %d\nNombre: %s\nLatitud: %.6f\nLongitud: %.6f", idGenerado, nombre, latitud, longitud))

		nombreEntry.SetText("")
		categoriaEntry.SetText("")
		ciudadEntry.SetText("")
		latitudEntry.SetText("")
		longitudEntry.SetText("")
	})

	longMinEntry := widget.NewEntry()
	longMinEntry.SetPlaceHolder("Ejemplo: -77.20")
	latMinEntry := widget.NewEntry()
	latMinEntry.SetPlaceHolder("Ejemplo: -12.30")
	longMaxEntry := widget.NewEntry()
	longMaxEntry.SetPlaceHolder("Ejemplo: -76.80")
	latMaxEntry := widget.NewEntry()
	latMaxEntry.SetPlaceHolder("Ejemplo: -11.90")

	formBuscar := widget.NewForm(
		widget.NewFormItem("Longitud mínima", longMinEntry),
		widget.NewFormItem("Latitud mínima", latMinEntry),
		widget.NewFormItem("Longitud máxima", longMaxEntry),
		widget.NewFormItem("Latitud máxima", latMaxEntry),
	)

	btnBuscar := widget.NewButton("Buscar por rango", func() { //rango
		minX, ok := leerFloatGUI(longMinEntry.Text, "longitud mínima", -180, 180, resultado)
		if !ok {
			return
		}
		minY, ok := leerFloatGUI(latMinEntry.Text, "latitud mínima", -90, 90, resultado)
		if !ok {
			return
		}
		maxX, ok := leerFloatGUI(longMaxEntry.Text, "longitud máxima", -180, 180, resultado)
		if !ok {
			return
		}
		maxY, ok := leerFloatGUI(latMaxEntry.Text, "latitud máxima", -90, 90, resultado)
		if !ok {
			return
		}

		if minX > maxX {
			minX, maxX = maxX, minX
		}
		if minY > maxY {
			minY, maxY = maxY, minY
		}

		zona := rtree.Rectangulo{MinX: minX, MinY: minY, MaxX: maxX, MaxY: maxY}
		lugares := arbol.BuscarPorRango(zona)
		resultado.SetText(formatearLugares("Resultados de búsqueda por rango", lugares))
	})

	idEliminarEntry := widget.NewEntry() //elimina
	idEliminarEntry.SetPlaceHolder("Ejemplo: 5")
	btnEliminar := widget.NewButton("Eliminar por ID", func() {
		id, err := strconv.Atoi(strings.TrimSpace(idEliminarEntry.Text))
		if err != nil {
			resultado.SetText("ID inválido. Debe ser un número entero.")
			return
		}

		eliminadoArbol := arbol.EliminarPorID(id)
		res, err := db.Exec("DELETE FROM lugares WHERE ID = @p1", id)
		if err != nil {
			resultado.SetText("Error al eliminar en SQL Server: " + err.Error())
			return
		}

		filas, _ := res.RowsAffected()
		if filas == 0 {
			resultado.SetText(fmt.Sprintf("No se encontró un lugar con ID %d.", id))
			return
		}

		if err := reiniciarIdentity(db); err != nil {
			resultado.SetText("Lugar eliminado, pero no se pudo reiniciar el IDENTITY: " + err.Error())
			return
		}

		if eliminadoArbol {
			resultado.SetText(fmt.Sprintf("Lugar con ID %d eliminado del R-Tree y de SQL Server.", id))
		} else {
			resultado.SetText(fmt.Sprintf("Lugar con ID %d eliminado de SQL Server. No estaba cargado en el R-Tree.", id))
		}
		idEliminarEntry.SetText("")
	})

	btnMostrarLugares := widget.NewButton("Mostrar lugares cargados", func() {
		lugares, err := cargarLugares(db)
		if err != nil {
			resultado.SetText("Error al cargar lugares: " + err.Error())
			return
		}
		resultado.SetText(formatearLugares("Lugares cargados en SQL Server", lugares))
	})

	btnMostrarArbol := widget.NewButton("Mostrar árbol", func() {
		resultado.SetText(arbol.TextoArbol())
	})

	btnRangoDemo := widget.NewButton("Rango demo Lima", func() {
		longMinEntry.SetText("-77.20")
		latMinEntry.SetText("-12.30")
		longMaxEntry.SetText("-76.80")
		latMaxEntry.SetText("-11.90")
		resultado.SetText("Se colocó un rango de demostración para Lima. Presiona 'Buscar por rango'.")
	})
	btnRangoGeneral := widget.NewButton("Mostrar rango general", func() {
		lugares, err := cargarLugares(db)
		if err != nil {
			resultado.SetText("Error al cargar lugares: " + err.Error())
			return
		}

		if len(lugares) == 0 {
			resultado.SetText("No hay lugares registrados para calcular el rango.")
			return
		}

		minLong, minLat, maxLong, maxLat := calcularRangoGeneral(lugares)

		longMinEntry.SetText(fmt.Sprintf("%.6f", minLong))
		latMinEntry.SetText(fmt.Sprintf("%.6f", minLat))
		longMaxEntry.SetText(fmt.Sprintf("%.6f", maxLong))
		latMaxEntry.SetText(fmt.Sprintf("%.6f", maxLat))

		resultado.SetText(fmt.Sprintf(
			"===== RANGO GENERAL DE LOS LUGARES =====\n\n"+
				"Longitud mínima: %.6f\n"+
				"Latitud mínima: %.6f\n"+
				"Longitud máxima: %.6f\n"+
				"Latitud máxima: %.6f\n",
			minLong,
			minLat,
			maxLong,
			maxLat,
		))
	})

	btnSalir := widget.NewButton("Salir", func() {
		ventana.Close()
	})

	panelIzquierdo := container.NewVScroll(container.NewVBox(
		widget.NewLabel("INSERTAR LUGAR"),
		formInsertar,
		btnInsertar,
		widget.NewSeparator(),
		widget.NewLabel("BUSCAR POR RANGO"),
		widget.NewLabel("Recuerde: LONGITUD (Este/Oeste [X]), LATITUD (Norte/Sur [Y])"),
		formBuscar,
		container.NewGridWithColumns(3, btnBuscar, btnRangoDemo, btnRangoGeneral), widget.NewSeparator(),
		widget.NewLabel("ELIMINAR"),
		widget.NewForm(widget.NewFormItem("ID", idEliminarEntry)),
		btnEliminar,
		widget.NewSeparator(),
		btnMostrarLugares,
		btnMostrarArbol,
		btnSalir,
	))

	contenido := container.NewHSplit(
		panelIzquierdo,
		container.NewBorder(widget.NewLabel("RESULTADOS"), nil, nil, nil, resultado),
	)
	contenido.Offset = 0.42

	ventana.SetContent(contenido)
	ventana.ShowAndRun()
}

func leerFloatGUI(texto string, campo string, minimo float64, maximo float64, resultado *widget.Entry) (float64, bool) {
	valor, err := strconv.ParseFloat(strings.TrimSpace(texto), 64)
	if err != nil {
		resultado.SetText(fmt.Sprintf("El campo %s debe ser un número decimal.", campo))
		return 0, false
	}

	if valor < minimo || valor > maximo {
		resultado.SetText(fmt.Sprintf("El campo %s debe estar entre %.0f y %.0f.", campo, minimo, maximo))
		return 0, false
	}

	return valor, true
}

func formatearLugares(titulo string, lugares []modelos.Lugar) string {
	var sb strings.Builder
	sb.WriteString("===== " + titulo + " =====\n")

	if len(lugares) == 0 {
		sb.WriteString("No hay lugares para mostrar.\n")
		return sb.String()
	}

	for _, lugar := range lugares {
		sb.WriteString(fmt.Sprintf(
			"ID: %d | Lugar: %s | Provincia/Dep: %s | Ciudad: %s | Lat: %.6f | Long: %.6f\n",
			lugar.Id,
			lugar.Nombre,
			lugar.Categoria,
			lugar.Ciudad,
			lugar.Latitud,
			lugar.Longitud,
		))
	}

	return sb.String()
}
func calcularRangoGeneral(lugares []modelos.Lugar) (float64, float64, float64, float64) {
	minLong := lugares[0].Longitud
	maxLong := lugares[0].Longitud
	minLat := lugares[0].Latitud
	maxLat := lugares[0].Latitud

	for _, lugar := range lugares {
		if lugar.Longitud < minLong {
			minLong = lugar.Longitud
		}

		if lugar.Longitud > maxLong {
			maxLong = lugar.Longitud
		}

		if lugar.Latitud < minLat {
			minLat = lugar.Latitud
		}

		if lugar.Latitud > maxLat {
			maxLat = lugar.Latitud
		}
	}

	return minLong, minLat, maxLong, maxLat
}
