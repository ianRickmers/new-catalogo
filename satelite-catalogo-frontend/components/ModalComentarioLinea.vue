<script setup>
import { ref, watch } from 'vue'

// Props
const props = defineProps({
  modelValue: Boolean,
  comentario: String,
  archivos: Array,
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

function guardar() {
  emit('guardar', comentarioLocal.value, archivosLocal.value)
  cerrar()
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
        />

        <!-- Campo archivos -->
        <v-file-input
          v-model="archivosLocal"
          label="Archivos adjuntos"
          multiple
          show-size
          prepend-icon="mdi-paperclip"

        />
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn text @click="cerrar">
          Cancelar
        </v-btn>

        <v-btn
          text
          color="primary"
          @click="guardar"
        >
          Guardar
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
