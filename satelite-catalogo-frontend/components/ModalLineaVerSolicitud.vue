<script setup>
import { ref, onMounted } from 'vue'
import ModalComentarioLinea from './ModalComentarioLinea.vue'
import { defineProps } from 'vue'
import { computed } from 'vue'
import ModalComentarioLineaVerSolicitud from './ModalComentarioLineaVerSolicitud.vue'


const props = defineProps({
  lineas: Array,
  productos: Array,
})
const mostrarModalComentario = ref(false)
const comentarioActual = ref('')
const archivosActuales = ref([])


// Combinar línea con su producto correspondiente
const lineasConProducto = computed(() => {
  return props.lineas.map((linea) => {
    const producto = props.productos.find(p => p?.id === linea.product_id)
    return {
      ...linea,
      descripcion: producto?.descripcion || '—',
      categoria: producto?.categoria || '—'
    }
  })
})

const headers = [
  { title: 'Línea', value: 'linea' },
  { title: 'Descripción', value: 'descripcion' },
  { title: 'Categoría', value: 'categoria' },
  { title: 'UM', value: 'um' },
  { title: 'Cantidad', value: 'cantidad' },
  { title: 'Precio', value: 'precio' },
  { title: 'Comentario', value: 'comentario' }
]

function abrirComentario(linea) {
  comentarioActual.value = linea.comentario || ''
  archivosActuales.value = linea.archivos || []
  mostrarModalComentario.value = true
}

</script>

<template>
  <!-- Encabezado -->
  <v-row
    style="background: linear-gradient(90deg, #00a499, #00b4a5); padding: 15px; border-radius: 5px; height: 50px; margin: 5px auto; justify-content: space-between;"
  >
    <h1 class="bold" style="color:#fff; font-size: 25px;">Líneas de la Solicitud</h1>
  </v-row>

  <!-- Tabla de líneas -->
  <div style="max-height: 300px; overflow-y: auto;">
  <v-data-table
    :headers="headers"
    :items="lineasConProducto"
    class="elevation-1 mt-4"
    hide-default-footer
    density="compact"
  >
    <template #item="{ item, index }">
      <tr>
        <td>{{ index + 1 }}</td>
        <td>{{ item.descripcion }}</td>
        <td>{{ item.categoria }}</td>
        <td>{{ item.um }}</td>
        <td>{{ item.cantidad }}</td>
        <td>{{ (item.importe_linea / item.cantidad) }}</td>
        <td>
          <v-btn icon color="white" @click="abrirComentario(item)">
            <v-icon>mdi-comment</v-icon>
          </v-btn>
        </td>
      </tr>
    </template>

    <template #no-data>
      <div class="text-center py-4">
        No hay líneas para mostrar.
      </div>
    </template>
  </v-data-table>
  </div>

  <!-- Modal para visualizar comentario -->
  <ModalComentarioLineaVerSolicitud
    v-model="mostrarModalComentario"
    :comentario="comentarioActual"
    :archivos="archivosActuales"
    solo-lectura
    @cerrar="mostrarModalComentario = false"
  />
</template>
