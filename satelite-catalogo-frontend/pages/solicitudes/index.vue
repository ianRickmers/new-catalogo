<script setup>
import axios from 'axios'
import { ref, computed, onMounted, watch, toRaw } from 'vue' // toRaw lo uso para evitar problemas con Proxy
import { useRouter } from 'vue-router'
import useAuthStore from '../store/useAuthStore'

definePageMeta({
  auth: true,
  userType: ['Administrador', 'Usuario']
})

const router = useRouter()
const authStore = useAuthStore()


// Estados para solicitudes y filtros
const allSolicitudes = ref([])
const idBusqueda = ref('')
const stateBusqueda = ref('')
const ccSeleccionado = ref(null)
const loading = ref(false)
const error = ref('')
const page = ref(1)
const pageSize = ref(50)
const totalPages = ref(0)
const menuFechas = ref(false)
const fechasSeleccionadas = ref([]) // máximo 2 fechas
// Mapas para mostrar número de centro de costo y etiquetas visuales
const ccMap = ref({})
const centrosDeCostosNumerados = ref([])

// Estados posibles de una solicitud
const estados = [
  { label: 'Todos', value: '' },
  { label: 'Inicial', value: 'I' },
  { label: 'Aprobada', value: 'A' },
  { label: 'Finalizada', value: 'C' },
  { label: 'Rechazada', value: 'D' },
  { label: 'Línea Aprobada', value: 'LA' },
  { label: 'Pendiente Aprobación', value: 'P' },
  { label: 'Revisión Preliminar', value: 'V' },
  { label: 'Cancelado', value: 'X' },
  { label: 'Abierta', value: 'O' }
]

// Computo los centros de costo disponibles para el selector
const centrosDeCostosDisponibles = computed(() => {
  return [{ label: 'Todos', value: null }, ...centrosDeCostosNumerados.value]
})

// Al cargar la vista, traigo centros de costo y luego solicitudes
onMounted(async () => {
  await cargarDatosCentrosCosto()
  buscarSolicitudes()

})

// Cada vez que cambia de página, actualizo resultados
watch(page, buscarSolicitudes)

function irACrearSolicitud() {
  router.push("/solicitudes/crearSolicitud")
}

// Solicito al backend las solicitudes paginadas y filtradas
async function buscarSolicitudes() {
  loading.value = true
  error.value = ''

  const params = {
    page: page.value,
    pageSize: pageSize.value
  }

  if (idBusqueda.value) params.id = idBusqueda.value
  if (stateBusqueda.value) params.state = stateBusqueda.value
  if (fechasSeleccionadas.value.length == 2) {
    // Si hay fechas seleccionadas, las agrego al filtro
    console.log("Entro a fechas")
    const [fecha1, fecha2] = [...fechasSeleccionadas.value].sort((a, b) => new Date(a) - new Date(b));
    params.fechaInicio = fecha1
    params.fechaFin = fecha2
  }
  // Obtener los CC del usuario autenticado
  const rawCCs = toRaw(authStore.getCC || [])
  const userCCs = Array.isArray(rawCCs) ? rawCCs : []

  // Si el usuario seleccionó un CC, lo uso; si no, uso todos los que tiene
  if (ccSeleccionado.value !== null) {
    params.ccs = ccSeleccionado.value
  } else if (userCCs.length > 0) {
    params.ccs = userCCs
  }

  try {
    const res = await axios.get(
      `${useRuntimeConfig().public.baseURL}/solicitud/filtradas`,
      {
        params,
        headers: {
          'Authorization': `Bearer ${authStore.getToken}`
        }
      }
    )

    allSolicitudes.value = res.data.data
    totalPages.value = res.data.totalPages || 1

    // Cargar visualmente los nombres/valores de CC para la tabla
    for (const s of allSolicitudes.value) {
      if (s.cc) cargarCentroCosto(s.cc)
    }
  } catch (err) {
    error.value = 'Error al filtrar solicitudes'
    allSolicitudes.value = []
    totalPages.value = 0
    console.error(err)
  } finally {
    loading.value = false
  }
}


// Consulto el número de un centro de costo por su ID
async function cargarCentroCosto(id) {
  if (ccMap.value[id]) return
  try {
    const res = await axios.get(`${useRuntimeConfig().public.baseURL}/cc/${id}`, {
      headers: {
        'Authorization': `Bearer ${authStore.getToken}`
      }
    })
    ccMap.value[id] = res.data.numero
  } catch (err) {
    ccMap.value[id] = '—'
  }
}


// Cargo todos los centros de costo del usuario desde authStore
async function cargarDatosCentrosCosto() {
  const rawIds = toRaw(authStore.getCC || [])
  const ids = Array.isArray(rawIds) ? rawIds : []
  console.log('Cargando centros de costo del usuario:', ids)

  const respuestas = await Promise.allSettled(
    ids.map(id => axios.get(`${useRuntimeConfig().public.baseURL}/cc/${id}`,{
      headers: {
        'Authorization': `Bearer ${authStore.getToken}`
      }
    }))
  )
  centrosDeCostosNumerados.value = respuestas
    .filter(r => r.status === 'fulfilled')
    .map(r => {
      const data = r.value.data
      return { label: `CC ${data.numero}`, value: data.id }
    })
}

// Cambiar filtros y refrescar
function buscarPorId() {
  page.value = 1
  buscarSolicitudes()
}

function onFiltroChange() {
  page.value = 1
  buscarSolicitudes()
}

function limpiarIdBusqueda() {
  idBusqueda.value = ''
  page.value = 1
  buscarSolicitudes()
}

function limpiarFiltros() {
  idBusqueda.value = ''
  stateBusqueda.value = ''
  ccSeleccionado.value = null
  page.value = 1
  fechasSeleccionadas.value = []
  buscarSolicitudes()
}

// Mostrar etiquetas y colores según el estado de solicitud
function formatStateLabel(state) {
  const map = {
    I: 'Inicial',
    A: 'Aprobada',
    C: 'Finalizada',
    D: 'Rechazada',
    LA: 'Línea Aprobada',
    P: 'Pendiente Aprobación',
    V: 'Revisión Preliminar',
    X: 'Cancelado',
    O: 'Abierta'
  }
  return map[state] || state
}

function mapColor(state) {
  switch (state) {
    case 'I': return 'grey-lighten-1'
    case 'A': return 'green'
    case 'C': return 'blue-grey'
    case 'D': return 'red'
    case 'LA': return 'light-green'
    case 'P': return 'amber'
    case 'V': return 'orange'
    case 'X': return 'deep-purple'
    case 'O': return 'indigo'
    default: return 'grey'
  }
}

function verSolicitud(id) {
  router.push(`/solicitudes/verSolicitud/${id}`)
}



// Mostrar texto en el input
const textoFechas = computed(() => {
  const opciones = { day: '2-digit', month: '2-digit', year: 'numeric' };

  if (fechasSeleccionadas.value.length === 2) {
    // Ordenar las fechas antes de mostrarlas
    const [fecha1, fecha2] = [...fechasSeleccionadas.value].sort((a, b) => new Date(a) - new Date(b));
    const fechaStr1 = new Date(fecha1).toLocaleDateString('es-CL', opciones);
    const fechaStr2 = new Date(fecha2).toLocaleDateString('es-CL', opciones);
    return `${fechaStr1} - ${fechaStr2}`;
  } else if (fechasSeleccionadas.value.length === 1) {
    return new Date(fechasSeleccionadas.value[0]).toLocaleDateString('es-CL', opciones);
  }

  return '';
});


// Controlar que no se seleccionen más de 2 fechas
watch(fechasSeleccionadas, (nuevas) => {
  if (nuevas.length > 2) {
    fechasSeleccionadas.value = nuevas.slice(-2)
  }

  // Si solo hay una fecha y se cerró el menú, limpiar
  if (nuevas.length === 1 && !menuFechas.value) {
    setTimeout(() => {
      if (fechasSeleccionadas.value.length === 1 && !menuFechas.value) {
        fechasSeleccionadas.value = []
      }
    }, 100) // espera a que el menú cierre completamente
  }
})


</script>



<template>
  <v-container fluid class="pa-0" style="height: 100vh;">
    <v-row no-gutters class="fill-height">
      <!-- Barra lateral -->
      <v-col cols="auto">
        <Navbar />
      </v-col>

      <!-- Contenido principal -->
      <v-col>
        <v-main class="bg-grey-lighten-5 pa-6">
          <v-container class="contenedor-solicitudes">

            <!-- Título y botón -->
            <v-row align="center" justify="space-between" class="mb-6">
              <h2 class="text-h5 font-weight-bold">Solicitudes</h2>
              <v-btn color="teal-darken-1" class="text-white" @click="irACrearSolicitud">
                + Nueva Solicitud
              </v-btn>
            </v-row>

            <!-- Filtros -->
            <v-row dense class="mb-4" align="center" justify="start">
              <v-col cols="12" md="3">
                <v-text-field
                  v-model="idBusqueda"
                  label="Buscar por ID"
                  variant="outlined"
                  density="compact"
                  @keyup.enter="buscarPorId"
                  clearable
                  @click:clear="limpiarIdBusqueda"
                />
              </v-col>

              <v-col cols="12" md="3">
                <v-select
                  v-model="stateBusqueda"
                  :items="estados"
                  item-title="label"
                  item-value="value"
                  label="Filtrar por estado"
                  variant="outlined"
                  density="compact"
                  @update:modelValue="onFiltroChange"
                />
              </v-col>

              <v-col cols="12" md="3" v-if="centrosDeCostosDisponibles.length">
                <v-select
                  v-model="ccSeleccionado"
                  :items="centrosDeCostosDisponibles"
                  item-title="label"
                  item-value="value"
                  label="Filtrar por Centro de Costo"
                  variant="outlined"
                  density="compact"
                  @update:modelValue="onFiltroChange"
                />
              </v-col>

              <v-col cols="12" md="3">
                <!-- Selector de fechas-->
                <v-menu
                  v-model="menuFechas"
                  :close-on-content-click="false"
                  transition="scale-transition"
                  offset-y
                  min-width="auto"  
                  persistent
                >
                  <template #activator="{ props }">
                    <v-text-field
                      v-bind="props"
                      label="Filtrar por intervalo de fechas"
                      v-model="textoFechas"
                      variant="outlined"
                      density="compact"
                      readonly
                      clearable
                      @click:clear="() => fechasSeleccionadas = []"
                    />
                  </template>

                  <v-card>
                    <v-date-picker
                      v-model="fechasSeleccionadas"
                      multiple
                      @update:model-value="(val) => {
                        if (val.length === 2) {
                          menuFechas = false;
                          onFiltroChange();
                        }
                      }"
                      title="Filtro de Fechas"
                    />
                  </v-card>
                </v-menu>
              </v-col>

              <v-col cols="12" class="text-end mt-2">
                <v-btn color="grey" variant="tonal" @click="limpiarFiltros">
                  Borrar Filtros
                </v-btn>
              </v-col>
            </v-row>

            <!-- Tabla de resultados -->
            <v-table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Descripción</th>
                  <th>Centro de Costo</th>
                  <th>Fecha</th>
                  <th>Estado</th>
                  <th>Acciones</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="!loading && allSolicitudes.length === 0">
                  <td colspan="5" class="text-center">
                    No hay resultados para los filtros aplicados.
                  </td>
                </tr>
                <tr v-for="solicitud in allSolicitudes" :key="solicitud.id">
                  <td>{{ solicitud.id }}</td>
                  <td>{{ solicitud.description }}</td>
                  <td>{{ ccMap[solicitud.cc] ?? '—' }}</td>
                  <td>{{ new Date(solicitud.fecha_solicitud).toISOString().slice(0, 10) }}</td>
                  <td>
                    <v-chip :color="mapColor(solicitud.state)" text-color="white" label>
                      {{ formatStateLabel(solicitud.state) }}
                    </v-chip>
                  </td>
                  <td>
                    <v-btn variant="text" color="teal" @click="verSolicitud(solicitud.id)">
                      Ver
                    </v-btn>
                  </td>
                </tr>
              </tbody>
            </v-table>

            <!-- Paginación y feedback -->
            <v-pagination
              v-if="totalPages > 1"
              v-model="page"
              :length="totalPages"
              class="mt-4"
            />

            <v-alert v-if="error" type="error" class="mt-4">{{ error }}</v-alert>
            <v-progress-linear v-if="loading" indeterminate color="teal" class="mt-2" />
          </v-container>
        </v-main>
      </v-col>
    </v-row>
  </v-container>
</template>
