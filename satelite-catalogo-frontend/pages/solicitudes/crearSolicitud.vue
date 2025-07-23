<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useNotifier } from '@/composables/useNotifier'
import axios from 'axios'
const { notify } = useNotifier()
definePageMeta({
  auth: true,
  userType: ['Administrador', 'Usuario']
})

const router = useRouter()

// Variables reactivas
const solicitantes = ref([])
const nombreSolicitud = ref('')
const description = ref('') 
const totalGeneral = ref(0)
const monedas = ref(['CLP', 'USD', 'EUR'])
const origenSeleccionado = ref({}) 
const cc = ref([])
const monedaSeleccionada = ref('CLP')
const fechaSolicitud = ref('')
const fechaContable = ref('')
const solicitanteSeleccionado = ref('')
const lineasRecibidas = ref([])
const archivosRecibidos = ref([]) // para almacenar archivos recibidos del componente hijo
const authStore = useAuthStore()
const centrosCostoUsuarioAutenticado = authStore.getCC || []
const mostrarConfirmacion = ref(false)

// Referencia a componente hijo, permite que el padre se comunique con el hijo
// para que este emita las lineas del formulario
const lineaRef = ref(null)


// Funciones para eventos emitidos por Linea
function actualizarTotalGlobal(total) {
  totalGeneral.value = total
}
// se recibe los archivos del componente hijo
function recibirArchivos(archivos) {
  archivosRecibidos.value = archivos
}
// se recibe las lineas del componente hijo
function recibirLineas(lineas) {
  lineasRecibidas.value = lineas
}
// se valida el formato de fecha YYYY-MM-DD
function esFechaValida(valor) {
  const regex = /^\d{4}-\d{2}-\d{2}$/
  if (!regex.test(valor)) return false

  const fecha = new Date(valor)
  return !isNaN(fecha.getTime())
}
// abre el diálogo de confirmación de envio de solicitud
function guardarFormulario() {
  mostrarConfirmacion.value = true
}
// validar campos del formulario antes de enviar
function validarFormulario() {
  let esValido = true

  if (!nombreSolicitud.value.trim()) {
    notify('El nombre de la solicitud es obligatorio.', 'error')
    esValido = false
  }

  if (!description.value.trim()) {
    notify('La descripción es obligatoria.', 'error')
    esValido = false
  }

  if (!fechaSolicitud.value) {
    notify('La fecha de solicitud es obligatoria.', 'error')
    esValido = false
  }

  if (!fechaContable.value) {
    notify('La fecha contable es obligatoria.', 'error')
    esValido = false
  }
  if (!esFechaValida(fechaSolicitud.value)) {
    notify('La fecha de solicitud tiene un formato incorrecto (YYYY-MM-DD).', 'error')
    esValido = false
  }

    if (!esFechaValida(fechaContable.value)){
    notify('La fecha contable tiene un formato incorrecto (YYYY-MM-DD).', 'error')
    esValido = false
  }
  if (!monedaSeleccionada.value) {
    notify('Debe seleccionar una moneda.', 'error')
    esValido = false
  }       


    if (!solicitanteSeleccionado.value) {
    notify('Debe seleccionar un solicitante.', 'error')
    esValido = false
  }

  if (!origenSeleccionado.value?.id) {
    notify('Debe seleccionar un centro de costo.', 'error')
    esValido = false
  }

  if (!lineasRecibidas.value.length) {
    notify('Debe agregar al menos una línea.', 'error')
    esValido = false
  }

  return esValido
}

// se lleva a cabo la confirmación del envío de la solicitud
function confirmarEnvio() {
  if (!process.client) return
  mostrarConfirmacion.value = false
  // Solicitar a Linea que emita las líneas actuales
  lineaRef.value.emitirLineasFormulario()

  // Usar nextTick en lugar de setTimeout para esperar a que lineas se actualicen
  nextTick(() => {
    // Validar campos requeridos
    if (!validarFormulario()) return

    enviarSolicitud()
    router.push({ path: '/solicitudes', query: { refresh: Date.now() } })
  })
}


// Cargar centros de costo disponibles
async function obtenerCC() {
  try {
    const response = await axios.get(`${useRuntimeConfig().public.baseURL}/cc/`, {
      headers: {
        'Authorization': `Bearer ${authStore.getToken}`
      }
    })

    const data = response.data
    cc.value = data.map(cc => ({
      nombre: cc.nombre,
      id: cc.id,
      numero: cc.numero,
      jefe: cc.jefe,
    }))

    // Seleccionar automáticamente el CC del usuario autenticado
    const centroCostoBuscado = cc.value.find(cc => centrosCostoUsuarioAutenticado.includes(cc.id))
    if (centroCostoBuscado) {
      origenSeleccionado.value = centroCostoBuscado
    }
  } catch (error) {
    if (error.response) {
      console.error('Error al cargar centro de costo:', error.response.data?.error || error.response.statusText)
    } else {
      console.error('Error de red al cargar centro de costo:', error.message)
    }
  }
}



// Obtener usuarios que pueden ser solicitantes
async function obtenerSolicitantesByCC(centrosCosto) {
  try {
    const response = await axios.post(
      `${useRuntimeConfig().public.baseURL}/user/by-cc`,
      centrosCosto,
      {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${authStore.getToken}`
        }
      }
    )

    const data = response.data
    solicitantes.value = data.map(user => ({
      nombre: user.username,
      id: user._id
    }))
  } catch (error) {
    if (error.response) {
      console.error('Error al cargar solicitantes:', error.response.data?.error || error.response.statusText)
    } else {
      console.error('Error de red al cargar solicitantes:', error.message)
    }
  }
}



// Enviar Solicitud al Backend
async function enviarSolicitud() {
  try {
    const formData = new FormData()

    const solicitud = {
      solicitante: solicitanteSeleccionado.value,
      fecha_solicitud: new Date(fechaSolicitud.value).toISOString(),
      fecha_contable: new Date(fechaContable.value).toISOString(),
      moneda: monedaSeleccionada.value,
      cc: origenSeleccionado.value.id,
      importe_total: totalGeneral.value,
      lines: lineasRecibidas.value,
      aprobador: origenSeleccionado.value.jefe,
      description: description.value,
      documents: [],
      state: "I",
      nombre_solicitud: nombreSolicitud.value
    }

    formData.append('solicitud', JSON.stringify(solicitud))

    archivosRecibidos.value.forEach((archivo) => {
      formData.append('archivos', archivo)
    })

    console.log('Datos de la solicitud:', formData.get('solicitud'))
    console.log('Archivos adjuntos:', formData.getAll('archivos'))

    const response = await axios.post(
      `${useRuntimeConfig().public.baseURL}/solicitud/`,
      formData,
      {
        headers: {
          'Authorization': `Bearer ${authStore.getToken}`
        }
      }
    )

    console.log('Solicitud enviada con éxito:', response.data)
    notify('Solicitud enviada con éxito', 'success')
  } catch (error) {
    const errorMsg = error.response?.data?.error || error.message
    console.error('Error al enviar la solicitud:', errorMsg)
    notify('Error al enviar la solicitud: ' + errorMsg, 'error')
  }
}


// Ejecutar al cargar página
onMounted(() => {
  obtenerSolicitantesByCC(centrosCostoUsuarioAutenticado)
  obtenerCC()
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
          <v-text-field label="Nombre Solicitud" v-model="nombreSolicitud" />
        </v-col>
        <v-col>
          <v-text-field label="Descripción" v-model="description" />
        </v-col>
      </v-row>

      <!-- Encabezado cabecera -->
      <v-row style="background: linear-gradient(90deg, #00a499, #00b4a5); padding: 15px; border-radius: 5px; height:50px; margin:5px auto;">
        <h1 class="bold" style="color:#fff; font-size: 25px;">Cabecera</h1>
      </v-row>

      <!-- Cabecera -->
      <v-row>
        <v-col>
          <v-autocomplete
            label="Solicitante"
            :items="solicitantes"
            v-model="solicitanteSeleccionado"
            item-title="nombre"
            item-value="id"
            dense
            outlined
            clearable
            placeholder="Buscar solicitante"
            hide-no-data
            autocomplete="off"
          />
          <v-text-field type="date" v-model="fechaSolicitud" label="Fecha Solicitud" outlined dense />
          <v-select
            :items="monedas"
            label="Código Moneda"
            v-model="monedaSeleccionada"
            outlined
            dense
            clearable
            placeholder="Selecciona una moneda"
          />
        </v-col>

        <v-col>
          <v-select
            :items="cc"
            item-title="nombre"
            :item-value="item => item"
            label="Origen"
            v-model="origenSeleccionado"
            outlined
            dense
            clearable
            placeholder="Selecciona un centro de costo"
          />
          <v-text-field type="date" v-model="fechaContable" label="Fecha Contable" outlined dense />
          <div style="background-color: #00a499; color:#fff; border-radius: 5px; height: 56px; display: flex; align-items: center; padding-left: 10px; font-size: 14px; margin-top: 8px;">
            <span>Importe Total: {{ totalGeneral }}</span>
          </div>
        </v-col>
      </v-row>

      <!-- Componente hijo -->
      <ModalLinea
        ref="lineaRef"
        @actualizarTotalGlobal="actualizarTotalGlobal"
        @recibirLineas="recibirLineas"
        @recibirArchivos="recibirArchivos"
      />

      <!-- Botones de acción -->
      <v-row class="justify-center mt-6" dense>
        <v-btn color="#00A499" class="ma-2" @click="guardarFormulario">
          <v-icon left>mdi-content-save</v-icon>
          Guardar
        </v-btn>
      </v-row>
        <v-dialog v-model="mostrarConfirmacion" max-width="500px">
    <v-card>
      <v-card-title class="text-h6">Confirmar envío</v-card-title>
      <v-card-text>¿Estás seguro de que deseas enviar esta solicitud?, al ser enviada no podrá editar nuevamente.</v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text @click="mostrarConfirmacion = false">Cancelar</v-btn>
        <v-btn color="primary" text @click="confirmarEnvio">Confirmar</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

    </v-form>
  </v-container>
  </v-col>
  </v-row>
  </v-container>
</template>
