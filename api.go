//go:build !fyne

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"proyecto_algoritmos/modelos"
	"proyecto_algoritmos/rtree"
)

type Servidor struct {
	db    *sql.DB
	arbol *rtree.ArbolR
	mu    sync.Mutex
}

type RangoRequest struct {
	MinX float64 `json:"minX"`
	MinY float64 `json:"minY"`
	MaxX float64 `json:"maxX"`
	MaxY float64 `json:"maxY"`
}

type ArbolResponse struct {
	Texto string `json:"texto"`
}

type RangoGeneralResponse struct {
	LongitudMinima float64 `json:"longitudMinima"`
	LatitudMinima  float64 `json:"latitudMinima"`
	LongitudMaxima float64 `json:"longitudMaxima"`
	LatitudMaxima  float64 `json:"latitudMaxima"`
}

func main() {
	fmt.Println("Iniciando backend Go para Vue.js...")

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error al abrir conexión:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error al conectar con SQL Server:", err)
	}

	arbol := rtree.NuevoArbol()

	lugares, err := cargarLugares(db)
	if err != nil {
		log.Fatal("Error al cargar lugares:", err)
	}

	for _, lugar := range lugares {
		arbol.Insertar(lugar)
	}

	servidor := &Servidor{
		db:    db,
		arbol: arbol,
	}

	http.HandleFunc("/api/lugares", servidor.handleLugares)
	http.HandleFunc("/api/buscar-rango", servidor.handleBuscarRango)
	http.HandleFunc("/api/arbol", servidor.handleArbol)
	http.HandleFunc("/api/rango-general", servidor.handleRangoGeneral)
	http.HandleFunc("/api/recargar", servidor.handleRecargar)

	fmt.Println("Backend listo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func habilitarCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return true
	}

	return false
}

func responderJSON(w http.ResponseWriter, estado int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(estado)
	json.NewEncoder(w).Encode(data)
}

func responderError(w http.ResponseWriter, estado int, mensaje string) {
	responderJSON(w, estado, map[string]string{
		"error": mensaje,
	})
}

func (s *Servidor) handleLugares(w http.ResponseWriter, r *http.Request) {
	if habilitarCORS(w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.listarLugares(w, r)

	case http.MethodPost:
		s.insertarLugarAPI(w, r)

	case http.MethodDelete:
		s.eliminarLugarAPI(w, r)

	default:
		responderError(w, http.StatusMethodNotAllowed, "Método no permitido.")
	}
}

func (s *Servidor) listarLugares(w http.ResponseWriter, r *http.Request) {
	lugares, err := cargarLugares(s.db)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al cargar lugares: "+err.Error())
		return
	}

	responderJSON(w, http.StatusOK, lugares)
}

func (s *Servidor) insertarLugarAPI(w http.ResponseWriter, r *http.Request) {
	var lugar modelos.Lugar

	if err := json.NewDecoder(r.Body).Decode(&lugar); err != nil {
		responderError(w, http.StatusBadRequest, "Datos inválidos.")
		return
	}

	lugar.Nombre = strings.TrimSpace(lugar.Nombre)
	lugar.Categoria = strings.TrimSpace(lugar.Categoria)
	lugar.Ciudad = strings.TrimSpace(lugar.Ciudad)

	if lugar.Nombre == "" || lugar.Categoria == "" || lugar.Ciudad == "" {
		responderError(w, http.StatusBadRequest, "Nombre, departamento/región y ciudad no pueden estar vacíos.")
		return
	}

	if lugar.Latitud < -90 || lugar.Latitud > 90 {
		responderError(w, http.StatusBadRequest, "La latitud debe estar entre -90 y 90.")
		return
	}

	if lugar.Longitud < -180 || lugar.Longitud > 180 {
		responderError(w, http.StatusBadRequest, "La longitud debe estar entre -180 y 180.")
		return
	}

	var idGenerado int

	err := s.db.QueryRow(
		"INSERT INTO lugares (Nombre, Categoria, Ciudad, Latitud, Longitud) OUTPUT INSERTED.ID VALUES (@p1, @p2, @p3, @p4, @p5)",
		lugar.Nombre,
		lugar.Categoria,
		lugar.Ciudad,
		lugar.Latitud,
		lugar.Longitud,
	).Scan(&idGenerado)

	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al insertar en SQL Server: "+err.Error())
		return
	}

	lugar.Id = idGenerado

	s.mu.Lock()
	s.arbol.Insertar(lugar)
	s.mu.Unlock()

	responderJSON(w, http.StatusCreated, lugar)
}

func (s *Servidor) eliminarLugarAPI(w http.ResponseWriter, r *http.Request) {
	idTexto := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idTexto)
	if err != nil {
		responderError(w, http.StatusBadRequest, "ID inválido.")
		return
	}

	s.mu.Lock()
	eliminadoArbol := s.arbol.EliminarPorID(id)
	s.mu.Unlock()

	resultado, err := s.db.Exec("DELETE FROM lugares WHERE ID = @p1", id)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al eliminar en SQL Server: "+err.Error())
		return
	}

	filas, _ := resultado.RowsAffected()

	if filas == 0 {
		responderError(w, http.StatusNotFound, fmt.Sprintf("No se encontró el lugar con ID %d.", id))
		return
	}

	if err := reiniciarIdentity(s.db); err != nil {
		responderError(w, http.StatusInternalServerError, "Se eliminó, pero no se pudo reiniciar el IDENTITY: "+err.Error())
		return
	}

	responderJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje":        fmt.Sprintf("Lugar con ID %d eliminado correctamente.", id),
		"eliminadoArbol": eliminadoArbol,
	})
}

func (s *Servidor) handleBuscarRango(w http.ResponseWriter, r *http.Request) {
	if habilitarCORS(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		responderError(w, http.StatusMethodNotAllowed, "Método no permitido.")
		return
	}

	var req RangoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responderError(w, http.StatusBadRequest, "Datos de rango inválidos.")
		return
	}

	if req.MinX < -180 || req.MinX > 180 || req.MaxX < -180 || req.MaxX > 180 {
		responderError(w, http.StatusBadRequest, "La longitud debe estar entre -180 y 180.")
		return
	}

	if req.MinY < -90 || req.MinY > 90 || req.MaxY < -90 || req.MaxY > 90 {
		responderError(w, http.StatusBadRequest, "La latitud debe estar entre -90 y 90.")
		return
	}

	if req.MinX > req.MaxX {
		req.MinX, req.MaxX = req.MaxX, req.MinX
	}

	if req.MinY > req.MaxY {
		req.MinY, req.MaxY = req.MaxY, req.MinY
	}

	zona := rtree.Rectangulo{
		MinX: req.MinX,
		MinY: req.MinY,
		MaxX: req.MaxX,
		MaxY: req.MaxY,
	}

	s.mu.Lock()
	resultados := s.arbol.BuscarPorRango(zona)
	s.mu.Unlock()

	responderJSON(w, http.StatusOK, resultados)
}

func (s *Servidor) handleArbol(w http.ResponseWriter, r *http.Request) {
	if habilitarCORS(w, r) {
		return
	}

	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "Método no permitido.")
		return
	}

	s.mu.Lock()
	texto := s.arbol.TextoArbol()
	s.mu.Unlock()

	responderJSON(w, http.StatusOK, ArbolResponse{
		Texto: texto,
	})
}

func (s *Servidor) handleRangoGeneral(w http.ResponseWriter, r *http.Request) {
	if habilitarCORS(w, r) {
		return
	}

	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "Método no permitido.")
		return
	}

	lugares, err := cargarLugares(s.db)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al cargar lugares: "+err.Error())
		return
	}

	if len(lugares) == 0 {
		responderError(w, http.StatusNotFound, "No hay lugares registrados.")
		return
	}

	minLong, minLat, maxLong, maxLat := calcularRangoGeneral(lugares)

	responderJSON(w, http.StatusOK, RangoGeneralResponse{
		LongitudMinima: minLong,
		LatitudMinima:  minLat,
		LongitudMaxima: maxLong,
		LatitudMaxima:  maxLat,
	})
}

func (s *Servidor) handleRecargar(w http.ResponseWriter, r *http.Request) {
	if habilitarCORS(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		responderError(w, http.StatusMethodNotAllowed, "Método no permitido.")
		return
	}

	lugares, err := cargarLugares(s.db)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al recargar lugares: "+err.Error())
		return
	}

	nuevoArbol := rtree.NuevoArbol()

	for _, lugar := range lugares {
		nuevoArbol.Insertar(lugar)
	}

	s.mu.Lock()
	s.arbol = nuevoArbol
	s.mu.Unlock()

	responderJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje":  "R-Tree recargado correctamente.",
		"cantidad": len(lugares),
	})
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