# Proyecto Algoritmos - R-Tree
Este proyecto implementa un R-Tree desde cero en Go para indexar lugares geográficos usando latitud y longitud.

## Integrantes
- Luanna Fiorella Fuentes Lopez
- Daniela Alexandra Rios Flores
- Emely Cristhel Asencios Gomez
- Dayana Koraly Jaimes Ostos

## Estructura usada
La estructura implementada es un R-Tree, una estructura espacial usada para organizar datos multidimensionales. En este proyecto se usa para almacenar lugares geográficos y realizar búsquedas por rango.

## Caso de uso
El sistema permite registrar lugares con:
- Nombre
- Provincia o departamento
- Ciudad
- Latitud
- Longitud

## El R-Tree usa:
- X = Longitud
- Y = Latitud

## Operaciones implementadas
- Inserción de lugares
- Búsqueda por rango
- Eliminación por ID
- Visualización del árbol
- Carga de datos desde SQL Server
- Simulación visual con Vue.js y backend en Go

## Ejecucion del frontend Vue

cd frontend-vue
npm.cmd install
npm.cmd run dev

http://localhost:5173

## Ejecución en bash
```bash
go run .
go mod tidy
go run .
go build .

## Backend 

http://localhost:8080

## unitarios
go test ./...

## Benchmarks
go test ./rtree -run='^$' -bench='.' -benchmem -count=1

## Complejidad Big-O
Inserción: O(log n) promedio.
Búsqueda por rango: O(log n + k) promedio, donde k es la cantidad de resultados encontrados
Eliminación: O(log n) promedio
Peor caso: O(n), cuando los rectángulos se solapan demasiado

## Base de datos
El proyecto usa SQL Server con una tabla llamada lugares.

## Campos usados:
ID
Nombre
Categoria
Ciudad
Latitud
Longitud
Simulación visual

## La interfaz permite demostrar visualmente:
Carga de lugares
Mostrar árbol
Buscar por rango
Insertar nuevos lugares
Eliminar lugares
Calcular rango general