package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"

	"proyecto_algoritmos/modelos"
	"proyecto_algoritmos/rtree"
)

const connString = "server=localhost;database=ProyectoRTree;user id=sa;password=123456789;encrypt=disable;connection timeout=5"

func mainTerminal() {
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("conexion exitosa con SQL Server")

	arbol := rtree.NuevoArbol()
	lugares, err := cargarLugares(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, lugar := range lugares {
		arbol.Insertar(lugar)
	}

	fmt.Println("se cargaron", len(lugares), "lugares desde la base de datos")

	lector := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		fmt.Println("===================================")
		fmt.Println("            MENÚ RTREE")
		fmt.Println("===================================")
		fmt.Println("1. Mostrar lugares cargados")
		fmt.Println("2. Mostrar árbol")
		fmt.Println("3. Buscar por rango")
		fmt.Println("4. Insertar lugar")
		fmt.Println("5. Eliminar lugar")
		fmt.Println("6. Salir")

		opcion := leerEntero(lector, "Seleccione una opción: ")

		switch opcion {
		case 1:
			lugares, err := cargarLugares(db)
			if err != nil {
				fmt.Println("Error al cargar lugares:", err)
				continue
			}
			mostrarLugares(lugares)

		case 2:
			arbol.MostrarArbol()

		case 3:
			buscarPorRango(lector, arbol)

		case 4:
			insertarLugar(lector, db, arbol)

		case 5:
			eliminarLugar(lector, db, arbol)

		case 6:
			fmt.Println("Hasta luego.")
			return

		default:
			fmt.Println("Opción inválida.")
		}
	}
}

func cargarLugares(db *sql.DB) ([]modelos.Lugar, error) {
	rows, err := db.Query("SELECT ID, Nombre, Categoria, Ciudad, Latitud, Longitud FROM lugares")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listaLugares []modelos.Lugar

	for rows.Next() {
		var lugar modelos.Lugar

		err := rows.Scan(
			&lugar.Id,
			&lugar.Nombre,
			&lugar.Categoria,
			&lugar.Ciudad,
			&lugar.Latitud,
			&lugar.Longitud,
		)
		if err != nil {
			return nil, err
		}

		listaLugares = append(listaLugares, lugar)
	}

	return listaLugares, rows.Err()
}

func mostrarLugares(lugares []modelos.Lugar) {
	fmt.Println("\n===== LUGARES CARGADOS =====")

	if len(lugares) == 0 {
		fmt.Println("No hay lugares registrados.")
		return
	}

	for _, lugar := range lugares {
		fmt.Println("--------------------------------")
		fmt.Println("ID:", lugar.Id)
		fmt.Println("Nombre:", lugar.Nombre)
		fmt.Println("Categoría:", lugar.Categoria)
		fmt.Println("Ciudad:", lugar.Ciudad)
		fmt.Println("Latitud:", lugar.Latitud)
		fmt.Println("Longitud:", lugar.Longitud)
	}
}

func buscarPorRango(lector *bufio.Reader, arbol *rtree.ArbolR) {
	fmt.Println("\nIngrese el rectángulo de búsqueda")
	fmt.Println("Recuerda: X = Longitud, Y = Latitud")

	minX := leerLongitud(lector, "Longitud mínima: ")
	minY := leerLatitud(lector, "Latitud mínima: ")
	maxX := leerLongitud(lector, "Longitud máxima: ")
	maxY := leerLatitud(lector, "Latitud máxima: ")

	if minX > maxX {
		minX, maxX = maxX, minX
	}

	if minY > maxY {
		minY, maxY = maxY, minY
	}

	zona := rtree.Rectangulo{
		MinX: minX,
		MinY: minY,
		MaxX: maxX,
		MaxY: maxY,
	}

	resultados := arbol.BuscarPorRango(zona)

	fmt.Println()

	if len(resultados) == 0 {
		fmt.Println("No se encontraron lugares en ese rango.")
		return
	}

	fmt.Println("Lugares encontrados:")
	for _, lugar := range resultados {
		fmt.Printf(
			"- ID: %d | %s | %s | %s | Lat: %.6f | Long: %.6f\n",
			lugar.Id,
			lugar.Nombre,
			lugar.Categoria,
			lugar.Ciudad,
			lugar.Latitud,
			lugar.Longitud,
		)
	}
}

func insertarLugar(lector *bufio.Reader, db *sql.DB, arbol *rtree.ArbolR) {
	fmt.Println("\n===== INSERTAR LUGAR =====")

	nombre := leerTexto(lector, "Nombre: ")
	categoria := leerTexto(lector, "Categoría: ")
	ciudad := leerTexto(lector, "Ciudad: ")
	latitud := leerLatitud(lector, "Latitud: ")
	longitud := leerLongitud(lector, "Longitud: ")

	var idGenerado int
	err := db.QueryRow(
		"INSERT INTO lugares (Nombre, Categoria, Ciudad, Latitud, Longitud) OUTPUT INSERTED.ID VALUES (@p1, @p2, @p3, @p4, @p5)",
		nombre,
		categoria,
		ciudad,
		latitud,
		longitud,
	).Scan(&idGenerado)

	if err != nil {
		fmt.Println("Error al insertar en SQL Server:", err)
		fmt.Println("Revisa si la columna se llama Categoria o [Categoría].")
		return
	}

	lugar := modelos.Lugar{
		Id:        idGenerado,
		Nombre:    nombre,
		Categoria: categoria,
		Ciudad:    ciudad,
		Latitud:   latitud,
		Longitud:  longitud,
	}

	arbol.Insertar(lugar)

	fmt.Println("Lugar insertado correctamente con ID:", idGenerado)
}

func eliminarLugar(lector *bufio.Reader, db *sql.DB, arbol *rtree.ArbolR) {
	fmt.Println("\n===== ELIMINAR LUGAR =====")

	id := leerEntero(lector, "Ingrese el ID a eliminar: ")

	eliminadoArbol := arbol.EliminarPorID(id)

	resultado, err := db.Exec("DELETE FROM lugares WHERE ID = @p1", id)
	if err != nil {
		fmt.Println("Error al eliminar en SQL Server:", err)
		return
	}

	filasAfectadas, _ := resultado.RowsAffected()

	if filasAfectadas > 0 {
		if err := reiniciarIdentity(db); err != nil {
			fmt.Println("El lugar se eliminó, pero no se pudo reiniciar el IDENTITY:", err)
			return
		}

		if eliminadoArbol {
			fmt.Println("Lugar eliminado correctamente del R-Tree y de SQL Server.")
			fmt.Println("El contador de ID fue actualizado correctamente.")
			return
		}

		fmt.Println("Se eliminó de SQL Server. No estaba en el R-Tree cargado.")
		fmt.Println("El contador de ID fue actualizado correctamente.")
		return
	}

	fmt.Println("No se encontró el ID.")
}

func reiniciarIdentity(db *sql.DB) error {
	_, err := db.Exec(`
		DECLARE @maxId INT;

		SELECT @maxId = ISNULL(MAX(ID), 0) FROM lugares;

		DBCC CHECKIDENT ('lugares', RESEED, @maxId);
	`)
	return err
}

func leerTexto(lector *bufio.Reader, mensaje string) string {
	for {
		fmt.Print(mensaje)
		texto, _ := lector.ReadString('\n')
		texto = strings.TrimSpace(texto)

		if texto != "" {
			return texto
		}

		fmt.Println("El texto no puede estar vacío.")
	}
}

func leerEntero(lector *bufio.Reader, mensaje string) int {
	for {
		fmt.Print(mensaje)
		texto, _ := lector.ReadString('\n')
		texto = strings.TrimSpace(texto)

		numero, err := strconv.Atoi(texto)
		if err == nil {
			return numero
		}

		fmt.Println("Dato inválido. Ingresa un número entero.")
	}
}

func leerFloat(lector *bufio.Reader, mensaje string) float64 {
	for {
		fmt.Print(mensaje)
		texto, _ := lector.ReadString('\n')
		texto = strings.TrimSpace(texto)

		numero, err := strconv.ParseFloat(texto, 64)
		if err == nil {
			return numero
		}

		fmt.Println("Dato inválido. Ingresa un número decimal.")
	}
}

func leerLatitud(lector *bufio.Reader, mensaje string) float64 {
	for {
		latitud := leerFloat(lector, mensaje)

		if latitud >= -90 && latitud <= 90 {
			return latitud
		}

		fmt.Println("Latitud inválida. Debe estar entre -90 y 90.")
	}
}

func leerLongitud(lector *bufio.Reader, mensaje string) float64 {
	for {
		longitud := leerFloat(lector, mensaje)

		if longitud >= -180 && longitud <= 180 {
			return longitud
		}

		fmt.Println("Longitud inválida. Debe estar entre -180 y 180.")
	}
}
