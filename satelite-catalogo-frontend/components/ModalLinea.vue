<script setup>
import { ref } from 'vue'
import ModalConsultaProducto from './ModalConsultaProducto.vue'
// constantes reactivas para manejar el estado de la pestaña y el modal
const tab = ref(0)
const mostrarModalComentario = ref(false)
const filaSeleccionada = ref(null)
const comentarioActual = ref('')
const archivosActuales = ref([])

const mostrarConsultaArtModal = ref(false)
const filas = ref([])
const productoSeleccionado = ref(null);
// Encabezados de la tabla
const headers = [
  { title: 'Línea', value: 'linea' },
  { title: 'Descripción', value: 'descripcion' },
  { title: 'Categoría', value: 'categoria' },
  { title: 'UM', value: 'um' }, 
  { title: 'Cantidad', value: 'cantidad' },
  { title: 'Precio', value: 'precio' },
  { title: 'Eliminar', value: 'eliminar'},
  { title: 'Comentario', value: 'comentario' },
]


// Función para agregar una nueva fila
function agregarFila(producto) {
  filas.value.push({
    descripcion: producto.descripcion || '',
    categoria: producto.categoria || '', 
    um: producto.UM || '', 
    cantidad: 1, // Cantidad por defecto
    precio: producto.precio || '',
    product_id: producto.id|| null, // aqui se le da el valor de id_product desde el producto mismo
    // linea recibe product_id no id_product
  })
  calcularTotalGlobal()
}

// Función para eliminar una fila y recalcular el total
function eliminarFila(index) {
  filas.value.splice(index, 1)
  calcularTotalGlobal()
}

// Funciones para el modal de comentarios
function abrirComentario(index) {
  filaSeleccionada.value = index
  comentarioActual.value = filas.value[index].comentario || ''
  archivosActuales.value = filas.value[index].archivos || []
  mostrarModalComentario.value = true
}

function guardarComentario(comentario, archivos) {
  if (filaSeleccionada.value !== null) {
    filas.value[filaSeleccionada.value].comentario = comentario
    filas.value[filaSeleccionada.value].archivos = archivos
  }
  mostrarModalComentario.value = false
  filaSeleccionada.value = null
  comentarioActual.value = ''
  archivosActuales.value = []
}

// Emisión de eventos al componente padre
const emitLinea = defineEmits(['actualizarTotalGlobal', 'recibirLineas', 'recibirArchivos'])

// Cálculo del total de la solicitud
function calcularTotalGlobal() {
  const total = filas.value.reduce((acc, fila) => {
    const cantidad = parseFloat(fila.cantidad) || 0
    const precio = parseFloat(fila.precio) || 0
    return acc + cantidad * precio
  }, 0)
  emitLinea('actualizarTotalGlobal', total)
}

// Emitir líneas al padre
function emitirLineasFormulario() {
  const archivosTotales = []  // Arreglo para almacenar todos los archivos adjuntos
  const lineas = filas.value.map((item, index) => {
  const numeroLinea = index + 1

  // Adjuntar archivos con nombre modificado
  if (item.archivos && item.archivos.length > 0) {
    const archivosConLinea = item.archivos.map((archivo) => {
      const nuevoNombre = `linea_${numeroLinea}__${archivo.name}`

      // Crear nuevo File con el nombre modificado
      return new File([archivo], nuevoNombre, {
        type: archivo.type,
        lastModified: archivo.lastModified,
      })
    })

    archivosTotales.push(...archivosConLinea)
  }

  return {
    numero_linea: numeroLinea,
    cantidad: parseInt(item.cantidad),
    importe_linea: parseFloat(item.precio) * parseInt(item.cantidad),
    product_id: item.product_id || null,
    comentario: item.comentario || '',
    um: item.um || '',
  }
})

  console.log('Emitiendo líneas del formulario:', lineas)
  console.log('Archivos totales a emitir:', archivosTotales)

  // Emitimos ambos datos si quieres
  emitLinea('recibirLineas', lineas)
  emitLinea('recibirArchivos', archivosTotales)  // Puedes emitirlo como un nuevo evento
}


function handleProductoSeleccionado(producto) {
  mostrarConsultaArtModal.value = false;
  if (producto) {
    productoSeleccionado.value = producto;
    console.log(producto.id);
    agregarFila(producto);
    console.log('Producto seleccionado:', producto);
  } else {
    console.log('No se seleccionó ningún producto');
  }
}

// la funcion defineEmits permite emitir eventos al componente padre
// para esto funciona la variable lineRef en el componente padre
defineExpose({
  emitirLineasFormulario
})
</script>

<template>
  <!-- Encabezado -->
  <v-row
    style="background: linear-gradient(90deg, #00a499, #00b4a5); padding: 15px; border-radius: 5px; height: 50px; margin: 5px auto; justify-content: space-between;"
  >
    <h1 class="bold" style="color:#fff; font-size: 25px;">Línea</h1>
    <v-btn color="white" @click="mostrarConsultaArtModal=true">
      <v-icon left>mdi-plus</v-icon> Agregar Articulo
    </v-btn>
  </v-row>

  <!-- Pestañas -->
  <v-tabs v-model="tab" bg-color="white" dark>
    <v-tab value="0">Detalles</v-tab>
  </v-tabs>

  <!-- Contenido del tab -->
  <v-window v-model="tab">
    <v-window-item value="0">
      <v-card flat>
        <v-card-text>
          <!-- Tabla dinámica -->
      <div style="max-height: 400px; overflow-y: auto;">
        <v-data-table
          :headers="headers"
          :items="filas"
          item-value="index"
          class="elevation-1"
          hide-default-footer
          density="compact"
          :items-per-page="-1"
        >
          <template #item="{ item, index }">
            <tr>
              <td class="py-2 text-center" style="width: 40px;">
                <span>{{ index + 1 }}</span>
              </td>

              <td class="py-2" style="min-width: 150px;">
                <span>{{ item.descripcion }}</span>
              </td>

              <td class="py-2" style="min-width: 100px;">
                <span>{{ item.categoria }}</span>
              </td>

              <td class="py-2 text-center" style="width: 80px;">
                <span>{{ item.um }}</span>
              </td>

              <td class="py-2 text-center" style="width: 80px;">
                <v-text-field
                  v-model="item.cantidad"
                  type="number"
                  dense
                  outlined
                  hide-details
                  min="1"
                  style="max-width: 70px;"
                  @input="calcularTotalGlobal"
                />
              </td>

              <td class="py-2 text-right pr-4" style="width: 100px;">
                <span>{{ item.precio }}</span>
              </td>

              <td class="py-2 text-center" style="width: 50px;">
                <v-btn icon color="white" @click="eliminarFila(index)">
                  <v-icon>mdi-delete</v-icon>
                </v-btn>
              </td>

              <td class="py-2 text-center" style="width: 50px;">
                <v-btn icon color="white" @click="abrirComentario(index)">
                  <v-icon>mdi-comment</v-icon>
                </v-btn>
              </td>
            </tr>
          </template>

          <template #no-data>
            <div class="text-center py-4">
              No hay artículos agregados.
            </div>
          </template>
        </v-data-table>



          <!-- Modal para comentarios -->
          <ModalComentarioLinea
            v-model="mostrarModalComentario"
            :comentario="comentarioActual"
            :archivos="archivosActuales"
            @cerrar="mostrarModalComentario = false"
            @guardar="guardarComentario"
          />
          <!-- Modal de consulta de artículos -->
        <ModalConsultaProducto
          v-model="mostrarConsultaArtModal"
          @cerrar="handleProductoSeleccionado"
        />
        </div>
        </v-card-text>
      </v-card>
    </v-window-item>
  </v-window>
</template>
