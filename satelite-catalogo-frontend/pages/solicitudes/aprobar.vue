<script setup>
import axios from 'axios'
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import useAuthStore from '../store/useAuthStore'

definePageMeta({
    auth: true,
    userType: ['Administrador']
})

const router = useRouter()
const authStore = useAuthStore()


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
const ccMap = ref({})
const centrosDeCostosNumerados = ref([])

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

const centrosDeCostosDisponibles = computed(() => {
    return [{ label: 'Todos', value: null }, ...centrosDeCostosNumerados.value]
})

onMounted(async () => {
    await cargarDatosCentrosCosto()
    obtenerSolicitudesParaAprobar()
})

watch(page, obtenerSolicitudesParaAprobar)

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

async function cargarCentroCosto(id) {
    if (ccMap.value[id]) return
    try {
        const res = await axios.get(`${useRuntimeConfig().public.baseURL}/cc/${id}`,{
            headers: {'Authorization': `Bearer ${authStore.getToken}` }
        })
        ccMap.value[id] = res.data.numero
    } catch (err) {
        ccMap.value[id] = '—'
    }
}

async function cargarDatosCentrosCosto() {
    const ids = authStore.getCC || []
    const respuestas = await Promise.allSettled(
        ids.map(id => axios.get(`${useRuntimeConfig().public.baseURL}/cc/${id}`),{
            headers: { 'Authorization': `Bearer ${authStore.getToken}` }
        })
    )
    centrosDeCostosNumerados.value = respuestas
        .filter(r => r.status === 'fulfilled')
        .map(r => {
            const data = r.value.data
            return { label: `CC ${data.numero}`, value: data.id }
        })
}

const userId = computed(() => authStore.getID)

async function obtenerSolicitudesParaAprobar() {
  loading.value = true
  error.value = ''

  const params = {
    page: page.value,
    pageSize: pageSize.value,
    state: stateBusqueda.value,
    id: idBusqueda.value,
    cc: ccSeleccionado.value,
    userId: userId.value // <-- se pasa al backend
  }
  if (fechasSeleccionadas.value.length == 2) {
    // Si hay fechas seleccionadas, las agrego al filtro
    console.log("Entro a fechas")
    const [fecha1, fecha2] = [...fechasSeleccionadas.value].sort((a, b) => new Date(a) - new Date(b));
    params.fechaInicio = fecha1
    params.fechaFin = fecha2
  }
  try {
    const res = await axios.get(
      `${useRuntimeConfig().public.baseURL}/solicitud/aprobar`,
      {
        params, // parámetros de la URL
        headers: {
          'Authorization': `Bearer ${authStore.getToken}`
        }
      }
    )
    allSolicitudes.value = res.data.data
    totalPages.value = res.data.totalPages || 1
    allSolicitudes.value.forEach(s => {
      if (s.cc) cargarCentroCosto(s.cc)
    })
  } catch (err) {
    error.value = 'Error al cargar solicitudes a aprobar'
    allSolicitudes.value = []
    totalPages.value = 0
    console.error(err)
  } finally {
    loading.value = false
  }
}



function limpiarIdBusqueda() {
    idBusqueda.value = ''
    page.value = 1
    obtenerSolicitudesParaAprobar()
}

function limpiarFiltros() {
  idBusqueda.value = ''
  stateBusqueda.value = ''
  ccSeleccionado.value = null
  page.value = 1
  fechasSeleccionadas.value = []
  obtenerSolicitudesParaAprobar()
}
function aprobarSolicitud(id) {
  axios.put(
    `${useRuntimeConfig().public.baseURL}/solicitud/${id}`,
    { state: 'A' }, 
    {
      headers: {
        'Authorization': `Bearer ${authStore.getToken}`
      }
    }
  )
    .then(() => obtenerSolicitudesParaAprobar())
    .catch(err => console.error('Error al aprobar:', err))
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
            <v-col cols="auto">
                <Navbar />
            </v-col>
            <v-col>
                <v-main class="bg-grey-lighten-5 pa-6">
                    <v-container class="contenedor-solicitudes">
                        <h2 class="text-h5 font-weight-bold mb-6">Solicitudes para Aprobar</h2>

                        <v-row dense class="mb-4" align="center">
                            <v-col cols="12" md="3">
                                <v-text-field v-model="idBusqueda" label="Buscar por ID" variant="outlined"
                                    density="compact" @keyup.enter="obtenerSolicitudesParaAprobar" clearable
                                    @click:clear="limpiarIdBusqueda" />
                            </v-col>

                            <v-col cols="12" md="3">
                                <v-select v-model="stateBusqueda" :items="estados" item-title="label" item-value="value"
                                    label="Filtrar por estado" variant="outlined" density="compact"
                                    @update:modelValue="obtenerSolicitudesParaAprobar" />
                            </v-col>

                            <v-col cols="12" md="3">
                                <v-select v-model="ccSeleccionado" :items="centrosDeCostosDisponibles"
                                    item-title="label" item-value="value" label="Filtrar por Centro de Costo"
                                    variant="outlined" density="compact"
                                    @update:modelValue="obtenerSolicitudesParaAprobar" />
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
                                        obtenerSolicitudesParaAprobar();
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
                        <!-- centrar los tds de la tabla -->
                        <v-table>
                        <thead>
                            <tr>
                            <th>ID</th>
                            <th>Descripción</th>
                            <th>Centro de Costo</th>
                            <th>Fecha</th>
                            <th>Estado</th>
                            <th>Acciones</th>
                            <th>Aprobación</th>
                            </tr>
                        </thead>
                        <tbody >
                            <tr v-if="!loading && allSolicitudes.length === 0">
                            <td colspan="5">No hay solicitudes pendientes de aprobación.</td>
                            </tr>
                            <tr v-for="solicitud in allSolicitudes" :key="solicitud.id">
                            <td>{{ solicitud.id }}</td>
                            <td>{{ solicitud.description }}</td>
                            <td>{{ ccMap[solicitud.cc] ?? '—' }}</td>
                            <td>{{ new Date(solicitud.fecha_solicitud).toISOString().slice(0,10)}}</td>
                            <td>
                                <v-chip :color="mapColor(solicitud.state)" text-color="white" label>
                                {{ formatStateLabel(solicitud.state) }}
                                </v-chip>
                            </td>
                            <td>
                                <v-btn variant="text" color="teal" @click="verSolicitud(solicitud.id)">Ver</v-btn>
                            </td>
                            <td>
                                <!-- <v-btn v-if="solicitud.state !== 'A'" color="green" @click="aprobarSolicitud(solicitud.id)" small>
                                Aprobar
                                </v-btn> -->
                                <v-btn :color="solicitud.state !== 'A' ? 'green' : 'grey'" :disabled="solicitud.state === 'A'"
                                @click="aprobarSolicitud(solicitud.id)" small>
                                Aprobar
                                </v-btn>
                            </td>
                            </tr>
                        </tbody>
                        </v-table>

                        <v-pagination v-if="totalPages > 1" v-model="page" :length="totalPages" class="mt-4" />

                        <v-alert v-if="error" type="error" class="mt-4">{{ error }}</v-alert>
                        <v-progress-linear v-if="loading" indeterminate color="teal" class="mt-2" />
                    </v-container>
                </v-main>
            </v-col>
        </v-row>
    </v-container>
</template>
