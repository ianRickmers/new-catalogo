<script setup>
import { ref, watch } from 'vue'

// Props
const props = defineProps({
  modelValue: Boolean,
  comentario: String,
  archivos: Array,
  soloLectura: Boolean 
})

const emit = defineEmits(['update:modelValue', 'guardar'])

const comentarioLocal = ref(props.comentario || '')
const archivosLocal = ref([])

watch(() => props.modelValue, (abierto) => {
  if (abierto) {
    comentarioLocal.value = props.comentario || ''
    archivosLocal.value = Array.isArray(props.archivos) ? [...props.archivos] : []
  }
})

function cerrar() {
  emit('update:modelValue', false)
}


// Extraer nombre legible del archivo desde su ruta
function extraerNombre(ruta) {
  const partes = ruta.split('/')
  return partes[partes.length - 1]
}
// Da la ruta para descargar los archivos
function getRutaCompleta(ruta) {
  const base = useRuntimeConfig().public.baseURL
  return `${base}${ruta}`
}
</script>

<template>
  <v-dialog
    :model-value="modelValue"
    @update:modelValue="$emit('update:modelValue', $event)"
    max-width="500"
  >
    <v-card>
      <v-card-title class="text-h6">Comentario del artículo</v-card-title>
      <v-card-text>
        <!-- Campo comentario -->
        <v-textarea
          v-model="comentarioLocal"
          label="Comentario"
          auto-grow
          readonly  
        />

        <!-- Campo archivos -->
        <div v-if="archivosLocal.length">
        <h4 class="mt-4 mb-2">Archivos adjuntos:</h4>
        <v-list dense>
            <v-list-item
            v-for="(ruta, index) in archivosLocal"
            :key="index"
            >
          <a
            :href="getRutaCompleta(ruta)"
            :download="extraerNombre(ruta)"
            target="_blank"
            rel="noopener noreferrer"
            style="color: #1976d2; text-decoration: underline; width: 100%;"
          >
                {{ extraerNombre(ruta) }}
            </a>
            </v-list-item>
        </v-list>
        </div>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn text @click="cerrar">
          Cerrar
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
