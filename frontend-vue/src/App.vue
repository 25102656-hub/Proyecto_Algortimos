<script setup>
import { ref, onMounted } from 'vue'

const API = 'http://localhost:8080/api'

const lugares = ref([])
const resultados = ref([])
const textoArbol = ref('')
const mensaje = ref('')

const nuevoLugar = ref({
  nombre: '',
  categoria: '',
  ciudad: '',
  latitud: '',
  longitud: ''
})

const rango = ref({
  minX: '',
  minY: '',
  maxX: '',
  maxY: ''
})

const idEliminar = ref('')

function mostrarMensaje(texto) {
  mensaje.value = texto
}

async function cargarLugares() {
  try {
    const res = await fetch(`${API}/lugares`)
    const data = await res.json()

    if (!res.ok) {
      mostrarMensaje(data.error || 'Error al cargar lugares.')
      return
    }

    lugares.value = data
    resultados.value = []
    textoArbol.value = ''
    mostrarMensaje(`Se cargaron ${data.length} lugares desde SQL Server.`)
  } catch (error) {
    mostrarMensaje('No se pudo conectar con el backend Go. Verifica que esté ejecutándose en http://localhost:8080')
  }
}

async function mostrarArbol() {
  try {
    const res = await fetch(`${API}/arbol`)
    const data = await res.json()

    if (!res.ok) {
      mostrarMensaje(data.error || 'Error al mostrar el árbol.')
      return
    }

    textoArbol.value = data.texto
    resultados.value = []
    mostrarMensaje('Estructura del R-Tree mostrada correctamente.')
  } catch (error) {
    mostrarMensaje('No se pudo conectar con el backend Go.')
  }
}

async function mostrarRangoGeneral() {
  try {
    const res = await fetch(`${API}/rango-general`)
    const data = await res.json()

    if (!res.ok) {
      mostrarMensaje(data.error || 'Error al calcular el rango general.')
      return
    }

    rango.value.minX = data.longitudMinima
    rango.value.minY = data.latitudMinima
    rango.value.maxX = data.longitudMaxima
    rango.value.maxY = data.latitudMaxima

    mostrarMensaje(
      `Rango general calculado correctamente:

Longitud mínima: ${data.longitudMinima}
Latitud mínima: ${data.latitudMinima}
Longitud máxima: ${data.longitudMaxima}
Latitud máxima: ${data.latitudMaxima}

Los valores ya fueron colocados en los campos de búsqueda.`
    )
  } catch (error) {
    mostrarMensaje('No se pudo conectar con el backend Go.')
  }
}

function rangoDemoLima() {
  rango.value.minX = -77.20
  rango.value.minY = -12.30
  rango.value.maxX = -76.80
  rango.value.maxY = -11.90

  textoArbol.value = ''
  resultados.value = []
  mostrarMensaje('Se colocó un rango de demostración para Lima. Ahora presiona "Buscar por rango".')
}

async function buscarPorRango() {
  try {
    const body = {
      minX: Number(rango.value.minX),
      minY: Number(rango.value.minY),
      maxX: Number(rango.value.maxX),
      maxY: Number(rango.value.maxY)
    }

    if (
      rango.value.minX === '' ||
      rango.value.minY === '' ||
      rango.value.maxX === '' ||
      rango.value.maxY === ''
    ) {
      mostrarMensaje('Debes completar todos los campos del rango.')
      return
    }

    const res = await fetch(`${API}/buscar-rango`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    })

    const data = await res.json()

    if (!res.ok) {
      mostrarMensaje(data.error || 'Error en la búsqueda por rango.')
      return
    }

    resultados.value = data
    textoArbol.value = ''

    if (data.length === 0) {
      mostrarMensaje('No se encontraron lugares dentro del rango ingresado.')
    } else {
      mostrarMensaje(`Se encontraron ${data.length} lugares dentro del rango.`)
    }
  } catch (error) {
    mostrarMensaje('No se pudo conectar con el backend Go.')
  }
}

async function insertarLugar() {
  try {
    if (
      nuevoLugar.value.nombre.trim() === '' ||
      nuevoLugar.value.categoria.trim() === '' ||
      nuevoLugar.value.ciudad.trim() === '' ||
      nuevoLugar.value.latitud === '' ||
      nuevoLugar.value.longitud === ''
    ) {
      mostrarMensaje('Debes completar todos los campos para insertar un lugar.')
      return
    }

    const body = {
      nombre: nuevoLugar.value.nombre,
      categoria: nuevoLugar.value.categoria,
      ciudad: nuevoLugar.value.ciudad,
      latitud: Number(nuevoLugar.value.latitud),
      longitud: Number(nuevoLugar.value.longitud)
    }

    const res = await fetch(`${API}/lugares`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    })

    const data = await res.json()

    if (!res.ok) {
      mostrarMensaje(data.error || 'Error al insertar lugar.')
      return
    }

    mostrarMensaje(`Lugar insertado correctamente con ID ${data.id}.`)

    nuevoLugar.value = {
      nombre: '',
      categoria: '',
      ciudad: '',
      latitud: '',
      longitud: ''
    }

    await cargarLugares()
  } catch (error) {
    mostrarMensaje('No se pudo conectar con el backend Go.')
  }
}

async function eliminarLugar() {
  try {
    if (idEliminar.value === '') {
      mostrarMensaje('Debes ingresar un ID para eliminar.')
      return
    }

    const res = await fetch(`${API}/lugares?id=${idEliminar.value}`, {
      method: 'DELETE'
    })

    const data = await res.json()

    if (!res.ok) {
      mostrarMensaje(data.error || 'Error al eliminar lugar.')
      return
    }

    mostrarMensaje(data.mensaje)
    idEliminar.value = ''

    await cargarLugares()
  } catch (error) {
    mostrarMensaje('No se pudo conectar con el backend Go.')
  }
}

async function recargarArbol() {
  try {
    const res = await fetch(`${API}/recargar`, {
      method: 'POST'
    })

    const data = await res.json()

    if (!res.ok) {
      mostrarMensaje(data.error || 'Error al recargar el R-Tree.')
      return
    }

    textoArbol.value = ''
    resultados.value = []
    mostrarMensaje(`${data.mensaje}\nCantidad: ${data.cantidad}`)
  } catch (error) {
    mostrarMensaje('No se pudo conectar con el backend Go.')
  }
}

onMounted(() => {
  cargarLugares()
})
</script>

<template>
  <div class="app">
    <header class="encabezado">
      <h1>Proyecto Algoritmos - R-Tree</h1>
      <p>
        Búsqueda espacial de lugares geográficos del Perú usando Go, SQL Server y Vue.js.
      </p>
    </header>

    <main class="layout">
      <section class="panel">
        <h2>Insertar lugar</h2>

        <label>Nombre</label>
        <input v-model="nuevoLugar.nombre" placeholder="Ejemplo: Parque Kennedy" />

        <label>Departamento / Región</label>
        <input v-model="nuevoLugar.categoria" placeholder="Ejemplo: Lima" />

        <label>Ciudad / Zona</label>
        <input v-model="nuevoLugar.ciudad" placeholder="Ejemplo: Miraflores" />

        <label>Latitud</label>
        <input v-model="nuevoLugar.latitud" placeholder="Ejemplo: -12.121120" />

        <label>Longitud</label>
        <input v-model="nuevoLugar.longitud" placeholder="Ejemplo: -77.029700" />

        <button class="principal" @click="insertarLugar">Insertar lugar</button>

        <hr />

        <h2>Buscar por rango</h2>
        <p class="nota">Recuerda: X = Longitud, Y = Latitud</p>

        <label>Longitud mínima</label>
        <input v-model="rango.minX" placeholder="Ejemplo: -77.20" />

        <label>Latitud mínima</label>
        <input v-model="rango.minY" placeholder="Ejemplo: -12.30" />

        <label>Longitud máxima</label>
        <input v-model="rango.maxX" placeholder="Ejemplo: -76.80" />

        <label>Latitud máxima</label>
        <input v-model="rango.maxY" placeholder="Ejemplo: -11.90" />

        <div class="botones">
          <button @click="buscarPorRango">Buscar por rango</button>
          <button @click="rangoDemoLima">Rango demo Lima</button>
          <button @click="mostrarRangoGeneral">Mostrar rango general</button>
        </div>

        <hr />

        <h2>Eliminar</h2>

        <label>ID a eliminar</label>
        <input v-model="idEliminar" placeholder="Ejemplo: 18" />

        <button class="peligro" @click="eliminarLugar">Eliminar por ID</button>

        <hr />

        <div class="botones">
          <button @click="cargarLugares">Mostrar lugares</button>
          <button @click="mostrarArbol">Mostrar árbol</button>
          <button @click="recargarArbol">Recargar R-Tree</button>
        </div>
      </section>

      <section class="resultados">
        <h2>Resultados</h2>

        <pre class="mensaje">{{ mensaje }}</pre>

        <div v-if="textoArbol" class="bloque">
          <h3>Estructura del árbol</h3>
          <pre class="arbol">{{ textoArbol }}</pre>
        </div>

        <div v-if="resultados.length > 0" class="bloque">
          <h3>Resultados de búsqueda</h3>

          <div class="tabla-contenedor">
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Lugar</th>
                  <th>Departamento</th>
                  <th>Ciudad / Zona</th>
                  <th>Latitud</th>
                  <th>Longitud</th>
                </tr>
              </thead>

              <tbody>
                <tr v-for="lugar in resultados" :key="lugar.id">
                  <td>{{ lugar.id }}</td>
                  <td>{{ lugar.nombre }}</td>
                  <td>{{ lugar.categoria }}</td>
                  <td>{{ lugar.ciudad }}</td>
                  <td>{{ lugar.latitud }}</td>
                  <td>{{ lugar.longitud }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div v-if="lugares.length > 0 && resultados.length === 0 && !textoArbol" class="bloque">
          <h3>Lugares cargados</h3>

          <div class="tabla-contenedor">
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Lugar</th>
                  <th>Departamento</th>
                  <th>Ciudad / Zona</th>
                  <th>Latitud</th>
                  <th>Longitud</th>
                </tr>
              </thead>

              <tbody>
                <tr v-for="lugar in lugares" :key="lugar.id">
                  <td>{{ lugar.id }}</td>
                  <td>{{ lugar.nombre }}</td>
                  <td>{{ lugar.categoria }}</td>
                  <td>{{ lugar.ciudad }}</td>
                  <td>{{ lugar.latitud }}</td>
                  <td>{{ lugar.longitud }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<style scoped>
:global(body) {
  margin: 0;
  min-width: 1200px;
  background: #f4f6f8;
}

:global(#app) {
  width: 100%;
  max-width: none;
}

.app {
  font-family: Arial, sans-serif;
  padding: 24px;
  background: #f4f6f8;
  color: #222;
  min-height: 100vh;
  box-sizing: border-box;
}

.encabezado {
  background: white;
  padding: 22px 28px;
  border-radius: 12px;
  margin-bottom: 24px;
  border: 1px solid #ddd;
}

.encabezado h1 {
  margin: 0 0 8px 0;
  font-size: 30px;
}

.encabezado p {
  margin: 0;
  color: #555;
  font-size: 16px;
}

.layout {
  display: grid;
  grid-template-columns: 420px minmax(800px, 1fr);
  gap: 24px;
  align-items: start;
}

.panel,
.resultados {
  background: white;
  padding: 24px;
  border-radius: 12px;
  border: 1px solid #ddd;
  box-sizing: border-box;
}

.panel {
  min-width: 420px;
}

.resultados {
  min-width: 800px;
}

h2 {
  margin-top: 0;
  font-size: 24px;
}

h3 {
  font-size: 20px;
  margin-bottom: 10px;
}

label {
  display: block;
  font-weight: bold;
  margin-top: 12px;
  margin-bottom: 4px;
}

input {
  width: 100%;
  padding: 10px;
  margin-top: 2px;
  box-sizing: border-box;
  border: 1px solid #bbb;
  border-radius: 7px;
  font-size: 15px;
}

button {
  margin-top: 12px;
  padding: 10px 14px;
  border: none;
  border-radius: 7px;
  cursor: pointer;
  background: #2d6cdf;
  color: white;
  font-weight: bold;
  font-size: 14px;
}

button:hover {
  opacity: 0.9;
}

.principal {
  background: #2d6cdf;
}

.peligro {
  background: #c0392b;
}

.botones {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.nota {
  font-size: 14px;
  color: #555;
  margin-top: -5px;
}

.mensaje {
  background: #f1f1f1;
  padding: 14px;
  border-radius: 8px;
  white-space: pre-wrap;
  font-size: 15px;
  line-height: 1.5;
  margin-bottom: 18px;
}

.arbol {
  background: #f1f1f1;
  padding: 14px;
  border-radius: 8px;
  white-space: pre;
  overflow-x: auto;
  overflow-y: auto;
  max-height: 650px;
  font-family: Consolas, 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
  box-sizing: border-box;
}

.bloque {
  margin-top: 18px;
}

.tabla-contenedor {
  width: 100%;
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 12px;
  font-size: 14px;
  min-width: 850px;
}

th,
td {
  border: 1px solid #ddd;
  padding: 9px;
  text-align: left;
}

th {
  background: #e9eef7;
  font-weight: bold;
}

tr:nth-child(even) {
  background: #fafafa;
}

hr {
  margin: 26px 0;
  border: none;
  border-top: 1px solid #ccc;
}
</style>