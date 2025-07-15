<script setup>
import { useRoute } from 'vue-router'
import { onMounted, ref, computed } from 'vue'
import axios from 'axios'
definePageMeta({
  auth: true,
  userType: ['Administrador', 'Usuario']
})
const authStore = useAuthStore()
// Leer ID desde la ruta
const route = useRoute()
const solicitudId = route.params.id

// variables reactivas
const cargando = ref(true)
const solicitud = ref(null)

const nombreSolicitud = ref('')
const description = ref('')
const totalGeneral = ref(0)
const monedaSeleccionada = ref('')
const fechaSolicitud = ref('')
const fechaContable = ref('')
const solicitanteSeleccionado = ref(null)
const origenSeleccionado = ref(null)
const lineasRecibidas = ref([]) // esto se le pasa al hijo ModalLineaVerSolicitud.vue
const productosDeLineas = ref([]) // esto se le pasa al hijo ModalLineaVerSolicitud.vue
const rutasArchivosAdjuntos = ref([]) 

// Variables para mostrar solo texto con v-model y que se actualicen bien
const solicitanteNombre = ref('')
const origenNombre = ref('')

// se obtiene la solicitud por ID
async function getSolicitudByID(id) {
  try {
    const response = await axios.get(`${useRuntimeConfig().public.baseURL}/solicitud/${id}`, {
      headers: {
        'Authorization': `Bearer ${authStore.getToken}`
      }
    })

    solicitud.value = response.data

    // Rellenar campos
    nombreSolicitud.value = solicitud.value.nombre_solicitud
    description.value = solicitud.value.description
    totalGeneral.value = solicitud.value.importe_total
    monedaSeleccionada.value = solicitud.value.moneda
    fechaSolicitud.value = solicitud.value.fecha_solicitud.slice(0, 10)
    fechaContable.value = solicitud.value.fecha_contable.slice(0, 10)
    lineasRecibidas.value = solicitud.value.lines
    rutasArchivosAdjuntos.value = solicitud.value.documents || []
  } catch (error) {
    if (error.response) {
      console.error('Error al cargar Solicitud By ID:', error.response.data?.error || error.response.statusText)
    } else {
      console.error('Error de red al cargar solicitud:', error.message)
    }
  }
}

// se obtiene el solicitante ligado a la solicitud
async function getSolicitanteByID(id) {
  try {
    const response = await axios.get(`${useRuntimeConfig().public.baseURL}/user/${id}`, {
      headers: {
        'Authorization': `Bearer ${authStore.getToken}`
      }
    })

    solicitanteSeleccionado.value = response.data
    console.log('Solicitante seleccionado:', solicitanteSeleccionado.value)
    solicitanteNombre.value = solicitanteSeleccionado.value?.username || ''
  } catch (error) {
    if (error.response) {
      console.error('Error al cargar Solicitante By ID:', error.response.data?.error || error.response.statusText)
    } else {
      console.error('Error de red al cargar solicitante:', error.message)
    }
  }
}

// se obtiene el centro de costo ligado a la solicitud
async function getCCByID(id) {
  try {
    const response = await axios.get(`${useRuntimeConfig().public.baseURL}/cc/${id}`, {
      headers: {
        'Authorization': `Bearer ${authStore.getToken}`
      }
    })

    origenSeleccionado.value = response.data
    origenNombre.value = origenSeleccionado.value?.nombre || ''
  } catch (error) {
    if (error.response) {
      console.error('Error al cargar CC By ID:', error.response.data?.error || error.response.statusText)
    } else {
      console.error('Error de red al cargar centro de costo:', error.message)
    }
  }
}

// se obtiene el producto por ID para mandarlo a las lineas
async function getProductByID(id) {
  try {
    const response = await axios.get(`${useRuntimeConfig().public.baseURL}/product/${id}`, {
      headers: {
        'Authorization': `Bearer ${authStore.getToken}`
      }
    })
    return response.data
  } catch (error) {
    if (error.response) {
      console.error('Error al cargar Producto By ID:', error.response.data?.error || error.response.statusText)
    } else {
      console.error('Error de red al cargar producto:', error.message)
    }
    return null
  }
}

const rutasPorLinea = computed(() => {
  const agrupadas = {}
  rutasArchivosAdjuntos.value.forEach(ruta => {
    const match = ruta.match(/linea_(\d+)\//) // se busca la linea a la que pertenece
    if (match) {
      const lineaId = parseInt(match[1]) // match [1] contiene el número de línea y se convierte a entero
      if (!agrupadas[lineaId]) {
        agrupadas[lineaId] = []
      }
      agrupadas[lineaId].push(ruta)
    }
  })

  return agrupadas
})

// En onMounted:
onMounted(async () => {

  await getSolicitudByID(solicitudId)
  console.log('Solicitud cargada:', lineasRecibidas.value)
  if (solicitud.value) {
    await getSolicitanteByID(solicitud.value.solicitante)
    await getCCByID(solicitud.value.cc)
  }

  // Obtener productos de las líneas
  productosDeLineas.value = await Promise.all(
    lineasRecibidas.value.map(async (linea) => {
      return await getProductByID(linea.product_id)
    })
  )
  console.log("rutas", rutasArchivosAdjuntos.value)
  // Asociar archivos a cada línea
  lineasRecibidas.value = lineasRecibidas.value.map((linea) => {
    return {
      ...linea,
      archivos: rutasPorLinea.value[linea.numero_linea] || []
    }
  })
  console.log('Líneas con productos:', lineasRecibidas.value)
})
</script>

<template>
<v-container fluid class="pa-0" style="height: 100vh;">
  <v-row no-gutters class="fill-height">
      <!-- Barra lateral -->
  <v-col cols="auto">
      <Navbar />
  </v-col>
  <v-col>
    <v-container class="pa-6 rounded-lg elevation-3" style="background-color: white; margin:50px auto;">
      <v-form>
        <!-- Encabezado -->
        <v-row style="background: linear-gradient(90deg, #00a499, #00b4a5); padding: 15px; border-radius: 5px; height:50px; margin:5px auto;">
          <h1 class="bold" style="color:#fff; font-size: 25px;">Nueva Solicitud</h1>
        </v-row>

        <!-- Datos principales -->
        <v-row>
          <v-col>
            <v-text-field label="Nombre Solicitud" v-model="nombreSolicitud" readonly />
          </v-col>
          <v-col>
            <v-text-field label="Descripción" v-model="description" readonly />
          </v-col>
        </v-row>

        <!-- Encabezado cabecera -->
        <v-row style="background: linear-gradient(90deg, #00a499, #00b4a5); padding: 15px; border-radius: 5px; height:50px; margin:5px auto;">
          <h1 class="bold" style="color:#fff; font-size: 25px;">Cabecera</h1>
        </v-row>

        <!-- Cabecera -->
        <v-row>
          <v-col>
            <!-- Solicitante como texto con v-model -->
            <v-text-field
              label="Solicitante"
              v-model="solicitanteNombre"
              readonly
            />
            <v-text-field
              type="date"
              v-model="fechaSolicitud"
              label="Fecha Solicitud"
              outlined
              dense
              readonly
            />
            <!-- Moneda solo texto -->
            <v-text-field
              label="Código Moneda"
              v-model="monedaSeleccionada"
              readonly
              outlined
              dense
            />
          </v-col>

          <v-col>
            <!-- Origen como texto con v-model -->
            <v-text-field
              label="Origen"
              v-model="origenNombre"
              readonly
              outlined
              dense
            />
            <v-text-field
              type="date"
              v-model="fechaContable"
              label="Fecha Contable"
              outlined
              dense
              readonly
            />
            <div
              style="background-color: #00a499; color:#fff; border-radius: 5px; height: 56px; display: flex; align-items: center; padding-left: 10px; font-size: 14px; margin-top: 8px;"
            >
              <span>Importe Total: {{ totalGeneral }}</span>
            </div>
          </v-col>
        </v-row>
              <!-- Componente hijo -->
      <ModalLineaVerSolicitud :lineas="lineasRecibidas" :productos="productosDeLineas"/>
      </v-form>
    </v-container>
  </v-col>
  </v-row>
</v-container>
</template>
